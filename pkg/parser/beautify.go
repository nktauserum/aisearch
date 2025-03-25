package parser

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	mdplug "github.com/JohannesKaufmann/html-to-markdown/plugin"
)

func HTMLtoMarkdown(html *string) (string, error) {
	converter := md.NewConverter("", true, nil)
	converter.Use(mdplug.GitHubFlavored())

	markdown, err := converter.ConvertString(*html)
	if err != nil {
		return "", err
	}

	return markdown, nil
}
