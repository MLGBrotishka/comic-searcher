package entity

import (
	"testing"
)

func TestComicRaw_GetData(t *testing.T) {
	comic := ComicRaw{
		ID:         123,
		URL:        "http://example.com",
		Transcript: "Text",
		Alt:        "Alt",
		Title:      "Title",
	}

	data := comic.GetData()

	expected := "TitleTextAlt"

	if data != expected {
		t.Errorf("Expected %s, but got %s", expected, data)
	}
}
