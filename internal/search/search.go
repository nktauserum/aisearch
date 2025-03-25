package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nktauserum/aisearch/config"
	"github.com/nktauserum/aisearch/shared"
)

func SearchTavily(query string) ([]shared.Source, error) {
	config := config.GetConfig()

	r := new(shared.TavilyRequest)
	r.Query = query + " -site:youtube.com -site:yandex.ru/video -site:vk.com"
	r.NumResult = 5
	r.Answer = false

	request, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", config.SearchEngine.URL, bytes.NewReader(request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.SearchEngine.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response shared.TavilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// Handle error
		return nil, err
	}

	return response.Results, nil

}

func SearchDuckduckgo(query string) ([]string, error) {
	// TODO: Убрать захардкоженое значение!!!
	var host = "http://127.0.0.1:8080/api/v1/search"

	r := new(shared.DucksearchRequest)
	r.Query = query + " -site:youtube.com -site:yandex.ru/video -site:vk.com"
	r.NumResult = 3

	request, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", host, bytes.NewReader(request))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response shared.DucksearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// Handle error
		return nil, err
	}

	return response.Result, nil
}
