package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	Port   int `json:"port"`
	OpenAI struct {
		URL         string  `json:"base_url"`
		Key         string  `json:"api_key"`
		Model       string  `json:"model"`
		Temperature float64 `json:"temperature"`
	} `json:"chat_api"`
	SearchEngine struct {
		URL string `json:"base_url"`
		Key string `json:"api_key"`
	} `json:"search_engine"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := instance.loadConfig(); err != nil {
			panic(err)
		}
	})
	return instance
}

func (c *Config) loadConfig() error {
	// TODO: сменить путь
	file, err := os.Open("/home/nktauserum/Документы/aisearch/config/config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}
