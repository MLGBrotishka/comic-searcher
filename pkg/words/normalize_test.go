package words

import (
	"my_app/pkg/xkcd"
	"testing"
)

func TestNormalizeComics(t *testing.T) {
	// Пример входных данных
	comics := []xkcd.Comic{
		{
			ID:         1,
			URL:        "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg",
			Transcript: "[[A boy sits in a barrel which is floating in an ocean.]]\nBoy: I wonder where I'll float next?\n[[The barrel drifts into the distance. Nothing else can be seen.]]\n{{Alt: Don't we all.}}",
			Alt:        "Don't we all.",
		},
		{
			ID:         2,
			URL:        "https://imgs.xkcd.com/comics/tree_cropped_(1).jpg",
			Transcript: "[[Two trees are growing on opposite sides of a sphere.]]\n{{Alt-title: 'Petit' being a reference to Le Petit Prince, which I only thought about halfway through the sketch}}",
			Alt:        "'Petit' being a reference to Le Petit Prince, which I only thought about halfway through the sketch",
		},
	}

	// Вызов функции NormalizeComics
	normalizedComics := NormalizeComics(comics, false)

	// Проверка результата
	if len(normalizedComics) != len(comics) {
		t.Errorf("Expected %d comics, got %d", len(comics), len(normalizedComics))
	}

	// Проверка значений
	for i := range comics {
		if normalizedComics[i].URL != comics[i].URL {
			t.Errorf("Expected URL %s, got %s", comics[i].URL, normalizedComics[i].URL)
		}

		if normalizedComics[i].ID != comics[i].ID {
			t.Errorf("Expected ID %d, got %d", comics[i].ID, normalizedComics[i].ID)
		}
	}

	t.Log(normalizedComics)
}
