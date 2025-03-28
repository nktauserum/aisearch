package answer

import (
	"encoding/json"

	"github.com/nktauserum/aisearch/pkg/ai/models"
)

func Fish(result chan string) {
	data := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
		"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}

	for _, fish := range data {
		response := models.Response{Content: fish}
		json_data, _ := json.Marshal(response)
		result <- string(json_data)
	}
}
