package downloader

import (
	"net/url"
	"path"
	"regexp"
	"strings"
)

// LinkManager отвечает за проверку домена и генерацию локальных путей
type LinkManager struct {
	base *url.URL
	// очищающий regexp для небезопасных символов в путях файловой системы
	safeRe *regexp.Regexp
}

// NewLinkManager создаёт LinkManager с базовым URL
func NewLinkManager(base *url.URL) *LinkManager {
	return &LinkManager{
		base:   base,
		safeRe: regexp.MustCompile(`[^a-zA-Z0-9\-\._/]+`),
	}
}

// IsInternal проверяет, принадлежит ли ссылка тому же хосту (без поддоменов).
// Здесь сравниваем Hostname() для простоты.
func (lm *LinkManager) IsInternal(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	// Если у ссылки нет хоста (относительная), считаем внутренней
	if u.Host == "" {
		return true
	}
	// Сравниваем hostname
	return strings.EqualFold(u.Hostname(), lm.base.Hostname())
}

// ToLocalPath переводит URL в локальный путь для сохранения (относительно OutputDir).
// Примеры:
//   - https://example.com/ -> example.com/index.html
//   - https://example.com/about -> example.com/about
//   - https://example.com/assets/img.png -> example.com/assets/img.png
//
// Если путь оканчивается на / — добавляем index.html.
// Query превращаем в безопасный суффикс.
func (lm *LinkManager) ToLocalPath(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		// на случай ошибки — вернуть что-то безопасное
		safe := lm.safeRe.ReplaceAllString(raw, "_")
		return safe
	}

	// используем hostname (если пустой — берём base host)
	host := u.Hostname()
	if host == "" {
		host = lm.base.Hostname()
	}

	p := u.Path
	if p == "" || strings.HasSuffix(p, "/") {
		p = path.Join(p, "index.html")
	}

	// если в конце путь не имеет расширения — для простоты оставляем как есть
	// но если путь оканчивается на '/', мы уже добавили index.html

	// query -> добавляем как суффикс к имени файла, заменив небезопасные символы
	if u.RawQuery != "" {
		// добавить _ и затем query, но сначала sanitize
		q := lm.safeRe.ReplaceAllString(u.RawQuery, "_")
		// если путь уже имеет расширение, добавим перед расширением
		ext := path.Ext(p)
		if ext != "" {
			base := strings.TrimSuffix(p, ext)
			p = base + "_" + q + ext
		} else {
			p = p + "_" + q
		}
	}

	// собрать итоговый путь: host + p
	full := path.Join(host, p)
	// очистить повторные слэши
	full = path.Clean(full)
	// удалить ведущий "./" если он есть
	full = strings.TrimPrefix(full, "./")
	// финальная санитаризация: заменить любые неподходящие символы
	full = lm.safeRe.ReplaceAllString(full, "_")

	return full
}
