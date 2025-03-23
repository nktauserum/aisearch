package answer

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gocolly/colly"

	"github.com/nktauserum/aisearch/internal/tools"
	"github.com/nktauserum/aisearch/shared"
)

func parse(url string) (shared.Website, error) {
	site := shared.Website{URL: url}

	c := colly.NewCollector(colly.AllowURLRevisit())
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Игнорировать проверку сертификатов
		},
	})

	var visibleText strings.Builder

	c.OnHTML("style, script", func(e *colly.HTMLElement) {
		e.DOM.Remove()
	})

	c.OnHTML("body h1, body h2, body h3, body p", func(e *colly.HTMLElement) {

		if e.DOM.Parents().Filter("nav, footer, header, aside, table, iframe").Length() == 0 {
			text := strings.TrimSpace(e.Text)
			if len(text) > 0 {
				visibleText.WriteString(text + "\n")
			}
		}
	})

	c.OnHTML("body a", func(e *colly.HTMLElement) {
		// Проверяем, что родительский элемент не является навигацией, футером или сайдбаром
		if e.DOM.Parents().Filter("nav, footer, header, aside, table, iframe").Length() == 0 {
			linkText := strings.TrimSpace(e.Text)
			if len(linkText) > 0 {
				visibleText.WriteString(linkText + " ")
			}
		}
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		site.Title = e.Text
	})

	err := c.Visit(site.URL)
	if err != nil {
		log.Println("visit error:", err)
		return shared.Website{}, err
	}

	site.Content = visibleText.String()

	return site, nil
}

func trafilatura(url string) (string, error) {
	// TODO: убрать значение пути
	cmd := exec.Command("/home/veliashev/.local/bin/trafilatura", "-u", url)
	output, err := cmd.Output()

	if err != nil {
		log.Println(string(output))
		return "", fmt.Errorf("trafilatura error: %v", err)
	}

	return string(output), nil
}

func GetContent(ctx context.Context, url string) (shared.Website, error) {
	content, err := tools.GetWebResourceContent(ctx, url)
	if err != nil {
		return shared.Website{}, err
	}

	ws := shared.Website{
		Content:  content.GetText(),
		URL:      url,
		Title:    content.GetTitle(),
		Sitename: content.GetSitename(),
	}

	return ws, nil
}
