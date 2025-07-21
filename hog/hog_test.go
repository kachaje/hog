package hog_test

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/hog/hog"
)

func TestDrawLine(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "face.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}
	bounds := img.Bounds()

	result := hog.DrawLine(image.Pt(bounds.Max.X/2, bounds.Max.Y/2), 0.5, 100, img, color.RGBA{R: 255})

	if result.Bounds() != bounds {
		t.Fatal("Test failed")
	}

	filename := "outputLine.png"

	outputFile, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		outputFile.Close()

		os.Remove(filename)
	}()

	err = jpeg.Encode(outputFile, result, nil)
	if err != nil {
		t.Fatal(err)
	}
}
