package parser

import (
	"context"
	"strings"

	"github.com/mrusme/journalist/crawler"
	"github.com/nktauserum/aisearch/shared"
	"go.uber.org/zap"
)

func ParseHTML(ctx context.Context, url string) (*shared.Website, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var crwlr *crawler.Crawler = crawler.New(logger)
	crwlr.SetLocation(url)
	crwlr.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// Используем cycleTLS для обхода защиты от ботов
	article, err := crwlr.GetReadable(true)
	if err != nil {
		return nil, err
	}

	var page shared.Website

	page.Title = article.Title
	page.URL = url
	page.HTML = removeEmptyLinks(article.ContentHtml)
	page.Sitename = article.SiteName

	return &page, nil
}

// removeEmptyLinks удаляет все ссылки, не содержащие текстовой информации
func removeEmptyLinks(html string) string {
	// Используем strings.Builder для эффективного построения строки
	var builder strings.Builder
	inLink := false
	inText := false
	for i := 0; i < len(html); i++ {
		if strings.HasPrefix(html[i:], "<a") {
			inLink = true
		}
		if inLink && strings.HasPrefix(html[i:], "</a>") {
			inLink = false
			if !inText {
				// Пропускаем пустую ссылку
				i += len("</a>") - 1
				continue
			}
			inText = false
		}
		if inLink && !inText && html[i] != ' ' && html[i] != '\n' {
			inText = true
		}
		builder.WriteByte(html[i])
	}
	return builder.String()
}

func GetContent(ctx context.Context, url string) (shared.Website, error) {

	content, err := ParseHTML(ctx, url)
	if err != nil {
		return shared.Website{}, err
	}

	// Удаляем все символы # для корректного отображения заголовков
	content.HTML = strings.ReplaceAll(content.HTML, "#", "")

	// Преобразуем HTML в Markdown
	content.Content, err = HTMLtoMarkdown(&content.HTML)
	if err != nil {
		return shared.Website{}, err
	}
	// Удаляем лишние переносы строк
	content.Content = strings.ReplaceAll(content.Content, "\n\n", "\n")

	return *content, nil
}
