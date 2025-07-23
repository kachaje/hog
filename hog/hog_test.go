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

func TestResize(t *testing.T) {
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

	newImg := h.ResizeShrink(img, 64, 128)

	grayImg := h.ImgToGray(newImg)

	filename := "outputShrink.jpg"

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

func TestGradX(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "flowerGray.jpg"))
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

	result := h.GradX(*grayImg, 32, 64)

	target := float32(0.073929965)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradX(*grayImg, 0, 0)

	target = float32(0.031128405)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradX(*grayImg, 190, 127)

	target = float32(0.0)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}

func TestGradY(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "flowerGray.jpg"))
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

	result := h.GradY(*grayImg, 32, 64)

	target := float32(0.03891051)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradY(*grayImg, 0, 0)

	target = float32(0.031128405)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}

	result = h.GradY(*grayImg, 190, 127)

	target = float32(0)

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}

func TestGradOrien(t *testing.T) {
	gx := float32(11)
	gy := float32(8)

	h := hog.HOG{}

	mag, ori, deg := h.GradOrien(gx, gy)

	targetMag := 13.601470508735444
	targetOri := 0.628796286415433
	targetDeg := 36.02737338510361

	if mag != targetMag {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", targetMag, mag)
	}

	if ori != targetOri {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", targetOri, ori)
	}

	if deg != targetDeg {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", targetDeg, deg)
	}
}
