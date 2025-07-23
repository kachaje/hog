package hog_test

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/hog/hog"
)

func TestImgToGray(t *testing.T) {
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

	grayImg := h.ImgToGray(img)

	filename := "outputGray.png"

	outputFile, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		outputFile.Close()

		os.Remove(filename)
	}()

	err = jpeg.Encode(outputFile, grayImg, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestResizeImg(t *testing.T) {
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

	newImg := h.ResizeImg(img, 64)

	filename := "outputResized.png"

	outputFile, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		outputFile.Close()

		os.Remove(filename)
	}()

	err = jpeg.Encode(outputFile, newImg, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGradX(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "thumbnailGray.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	h := hog.HOG{}

	result := h.GradX(img, 95, 64)

	target := float32(0.04669261)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradX(img, 0, 0)

	target = float32(0.027237354)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradX(img, 190, 127)

	target = float32(-0.027237354)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}

func TestGradY(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "thumbnailGray.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	h := hog.HOG{}

	result := h.GradY(img, 95, 64)

	target := float32(-0.027237356)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradY(img, 0, 0)

	target = float32(0.027237354)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradY(img, 190, 127)

	target = float32(-0.027237354)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}
