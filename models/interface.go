package models

type ResponseSet struct {
	Query   string            `json:"query"`
	Include string            `json:"include"`
	Options map[string]string `json:"options"`
	Result  []string          `json:"result"`
}

type RequestSet struct {
	Query   string
	Exclude string
	Not     map[int]string
}

type SearcherSet struct {
	Include  string
	Searcher map[int][]rune
}
