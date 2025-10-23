package main

import (
	"flag"
	"fmt"
	"log"
	"ntptime/L2.16/internal/downloader"
	"os"
)

var (
	url           = flag.String("url", "", "URL для скачивания")
	depth         = flag.Int("depth", 1, "Глубина рекурсии")
	outputDir     = flag.String("output", "./result", "Выходная директория")
	maxConcurrent = flag.Int("concurrent", 5, "Максимум одновременных загрузок")
)

func main() {
	flag.Parse()

	if *url == "" {
		fmt.Println("Ошибка: необходимо указать URL")
		fmt.Println("Использование:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := &downloader.Config{
		BaseURL:       *url,
		Depth:         *depth,
		OutputDir:     *outputDir,
		MaxConcurrent: *maxConcurrent,
	}

	dl, err := downloader.New(config)
	if err != nil {
		log.Fatalf("Ошибка создания загрузчика: %v", err)
	}

	fmt.Printf("Начинаем загрузку %s в %s (глубина: %d)\n", *url, *outputDir, *depth)

	err = dl.Start()
	if err != nil {
		log.Fatalf("Ошибка загрузки: %v", err)
	}

	fmt.Println("Загрузка завершена успешно!")

}
