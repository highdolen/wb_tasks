package downloader

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Parser — парсер HTML для извлечения ссылок
type Parser struct{}

// NewParser создает экземпляр Parser
func NewParser() *Parser {
	return &Parser{}
}

// IsHTML — быстрая эвристика: есть ли в содержимом HTML
// (проверяем наличие <html или <!doctype или <body)
func (p *Parser) IsHTML(content []byte) bool {
	lc := bytes.ToLower(content)
	return bytes.Contains(lc, []byte("<html")) ||
		bytes.Contains(lc, []byte("<!doctype")) ||
		bytes.Contains(lc, []byte("<body"))
}

// ExtractLinks извлекает абсолютные ссылки из HTML-контента.
// baseURL используется для разрешения относительных ссылок.
func (p *Parser) ExtractLinks(content []byte, baseURL string) []string {
	out := make([]string, 0)
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		// Если не парсится как HTML, вернём пустой список.
		return out
	}

	base, _ := url.Parse(baseURL)

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, a := range n.Attr {
					if a.Key == "href" {
						if u := p.resolve(base, a.Val); u != "" {
							out = append(out, u)
						}
					}
				}
			case "img", "script", "source":
				for _, a := range n.Attr {
					if a.Key == "src" {
						if u := p.resolve(base, a.Val); u != "" {
							out = append(out, u)
						}
					}
					// обработка srcset
					if a.Key == "srcset" {
						surls := p.parseSrcset(a.Val)
						for _, s := range surls {
							if u := p.resolve(base, s); u != "" {
								out = append(out, u)
							}
						}
					}
				}
			case "link":
				var isStylesheet bool
				for _, a := range n.Attr {
					if a.Key == "rel" && strings.Contains(strings.ToLower(a.Val), "stylesheet") {
						isStylesheet = true
						break
					}
				}
				if isStylesheet {
					for _, a := range n.Attr {
						if a.Key == "href" {
							if u := p.resolve(base, a.Val); u != "" {
								out = append(out, u)
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	// Убираем дубликаты (цифровая простая фильтрация)
	seen := make(map[string]struct{}, len(out))
	result := make([]string, 0, len(out))
	for _, u := range out {
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		result = append(result, u)
	}
	return result
}

// resolve: превратить относительную ссылку в абсолютную на основе base.
// возвращает "" для невалидных href (javascript:, mailto: и т.п.)
func (p *Parser) resolve(base *url.URL, href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}
	// игнорируем javascript:, mailto:, data: (если хотим можно включать data: для inline)
	lower := strings.ToLower(href)
	if strings.HasPrefix(lower, "javascript:") || strings.HasPrefix(lower, "mailto:") || strings.HasPrefix(lower, "data:") {
		return ""
	}
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	res := base.ResolveReference(u)
	// убрать фрагмент (#...)
	res.Fragment = ""
	return res.String()
}

// parseSrcset разбирает атрибут srcset и возвращает список URL (без размеров)
func (p *Parser) parseSrcset(val string) []string {
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		// каждый part может быть "url 2x" или "url 300w"
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		fields := strings.Fields(p)
		if len(fields) > 0 {
			out = append(out, fields[0])
		}
	}
	return out
}
