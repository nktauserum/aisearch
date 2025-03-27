package shared

type SearchRequest struct {
	Query      string `json:"query"`
	FormatMode string `json:"format_mode,omitempty"`
}

type SearchResponse struct {
	Session  SearchSession `json:"session"`
	Response Research      `json:"response"`
}

type RefineResponse struct {
	Session  SearchSession `json:"session"`
	Response string        `json:"response"`
}

type Research struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type RefineRequest struct {
	UUID  string `json:"uuid"`
	Query string `json:"query"`
}

// TODO: Когда буду добавлять авторизацию, добавить сюда поле с userid
type SearchSession struct {
	UUID  string `json:"uuid"`
	Topic string `json:"topic"`
}

type Website struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	HTML     string `json:"html"`
	Sitename string `json:"sitename"`
}
