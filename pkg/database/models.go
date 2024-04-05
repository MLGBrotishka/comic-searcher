package db

type ComicsMap map[int]Comic

type Comic struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}
