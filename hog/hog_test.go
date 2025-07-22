package hog_test

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/hog/hog"
	_ "image/jpeg"
)

func TestCalculateGradient(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "flower.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	h := hog.HOG{}

	template := []any{-1, 0, 1}

	result := h.CalculateGradient(img, template)

	_ = result
}
