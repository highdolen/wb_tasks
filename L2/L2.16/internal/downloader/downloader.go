package downloader

import (
	"fmt"
	"log"
	"net/url"
	"sync"
)

// Task — задание для воркера
type Task struct {
	URL   string
	Depth int
}

// Downloader — главный объект загрузчика
type Downloader struct {
	config      *Config
	storage     *Storage
	fetcher     *Fetcher
	parser      *Parser
	linkManager *LinkManager

	visitedMu   sync.Mutex
	visitedURLs map[string]bool
	wg          sync.WaitGroup
}

// New — конструктор
func New(config *Config) (*Downloader, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("некорректный базовый URL '%s': %v", config.BaseURL, err)
	}

	return &Downloader{
		config:      config,
		storage:     NewStorage(config.OutputDir),
		fetcher:     NewFetcher(),
		parser:      NewParser(),
		linkManager: NewLinkManager(baseURL),
		visitedURLs: make(map[string]bool),
	}, nil
}

// Start — точка входа, запускает воркеров
func (d *Downloader) Start() error {
	log.Printf("Начинаем загрузку: %s (глубина: %d, воркеров: %d)\n",
		d.config.BaseURL, d.config.Depth, d.config.MaxConcurrent)

	jobs := make(chan Task, 100)

	// Запускаем воркеров
	for i := 0; i < d.config.MaxConcurrent; i++ {
		d.wg.Add(1)
		go d.worker(i+1, jobs)
	}

	// Отправляем первую задачу
	jobs <- Task{URL: d.config.BaseURL, Depth: 0}

	// Ждём, пока все воркеры закончат
	d.wg.Wait()
	close(jobs)

	log.Println("Загрузка завершена.")
	return nil
}

// worker — функция воркера
func (d *Downloader) worker(id int, jobs chan Task) {
	defer d.wg.Done()

	for task := range jobs {
		if d.checkAndMark(task.URL) {
			continue
		}

		log.Printf("[worker %d] Скачиваю: %s (глубина %d)\n", id, task.URL, task.Depth)
		content, err := d.fetcher.Fetch(task.URL)
		if err != nil {
			log.Printf("[worker %d] Ошибка скачивания %s: %v", id, task.URL, err)
			continue
		}

		// Сохраняем страницу/ресурс
		localPath := d.linkManager.ToLocalPath(task.URL)
		if err := d.storage.Save(localPath, content); err != nil {
			log.Printf("[worker %d] Ошибка сохранения %s: %v", id, task.URL, err)
			continue
		}

		// Если это HTML и глубина не превышена — ищем ссылки
		if task.Depth < d.config.Depth && d.parser.IsHTML(content) {
			links := d.parser.ExtractLinks(content, task.URL)
			for _, link := range links {
				if d.linkManager.IsInternal(link) && !d.isVisited(link) {
					select {
					case jobs <- Task{URL: link, Depth: task.Depth + 1}:
					default:
						// Если буфер заполнен — просто пропускаем
					}
				}
			}
		}
	}
}

// checkAndMark — потокобезопасная отметка, что URL уже скачан
func (d *Downloader) checkAndMark(u string) bool {
	d.visitedMu.Lock()
	defer d.visitedMu.Unlock()

	if d.visitedURLs[u] {
		return true
	}
	d.visitedURLs[u] = true
	return false
}

// isVisited — просто проверка без отметки
func (d *Downloader) isVisited(u string) bool {
	d.visitedMu.Lock()
	defer d.visitedMu.Unlock()
	return d.visitedURLs[u]
}
