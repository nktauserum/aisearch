package parser

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrusme/journalist/crawler"
	"github.com/nktauserum/aisearch/shared"
	"go.uber.org/zap"
)

func ParseHTML(ctx context.Context, url string) (*shared.Website, error) {
	select {
	case <-ctx.Done():
		fmt.Printf("context deadline: %s\n", url)
		return nil, nil
	default:
	}

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
	//	page.Content = article.ContentText
	page.Sitename = article.SiteName

	const maxLength = 20000
	if len(article.ContentText) > maxLength {
		page.Content = article.ContentText[:maxLength]
	} else {
		page.Content = article.ContentText
	}

	return &page, nil
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
