package entity

type ComicRaw struct {
	ID         int    `json:"num"`
	URL        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Title      string `json:"title"`
}

func (c ComicRaw) GetData() string {
	return c.Title + c.Transcript + c.Alt
}

type Comic struct {
	ID       int
	URL      string
	Keywords KeywordMap
}

type IdMap map[int]bool

type KeywordMap map[string]bool
