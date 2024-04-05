package xkcd

type Comic struct {
	ID         int    `json:"num"`
	URL        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}
