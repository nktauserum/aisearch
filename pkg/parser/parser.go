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
	//page.HTML = article.ContentHtml
	page.Content = article.ContentText
	page.Sitename = article.SiteName

	return &page, nil
}

func processLinks(html string) string {
	var builder strings.Builder
	inLink := false
	inText := false

	for i := 0; i < len(html); i++ {
		// Проверяем, открывается ли ссылка
		if strings.HasPrefix(html[i:], "<a") {
			inLink = true
			// Находим конец тега <a>
			for html[i] != '>' {
				i++
			}
			i++ // Пропускаем '>'
			continue
		}

		// Проверяем, закрывается ли ссылка
		if inLink && strings.HasPrefix(html[i:], "</a>") {
			inLink = false
			// Пропускаем закрывающий тег </a>
			for html[i] != '>' {
				i++
			}
			continue
		}

		// Если мы находимся внутри ссылки, проверяем, есть ли текст
		if inLink && html[i] != ' ' && html[i] != '\n' {
			inText = true
		}

		// Если текст внутри ссылки найден, записываем его
		if inText {
			builder.WriteByte(html[i])
		} else if !inLink {
			// Записываем обычный текст вне ссылок
			builder.WriteByte(html[i])
		}
	}

	return builder.String()
}

func GetContent(ctx context.Context, url string) (shared.Website, error) {
	select {
	case <-ctx.Done():
		return shared.Website{}, ctx.Err()
	default:
	}

	content, err := ParseHTML(ctx, url)
	if err != nil {
		return shared.Website{}, err
	}

	content.Content = strings.ReplaceAll(content.Content, "\n\n", "\n")

	return *content, nil
}
