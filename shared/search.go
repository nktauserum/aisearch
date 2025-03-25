package shared

type DucksearchRequest struct {
	Query     string `json:"query"`
	NumResult int    `json:"num_results"`
}

type DucksearchResponse struct {
	Result []string `json:"result"`
}

type TavilyRequest struct {
	Query     string `json:"query"`
	NumResult int    `json:"max_results"`
	Answer    bool   `json:"include_answer"`
}

type TavilyResponse struct {
	Results []Source `json:"results"`
}
