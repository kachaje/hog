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

func TestAngleWeight(t *testing.T) {
	h := hog.HOG{}

	s, p1, e, p2 := h.AngleWeight(13.6, 36)

	if s != 1 {
		t.Fatalf("Test failed. Expected: 1; Actual: %v", s)
	}

	if p1 != 2.72 {
		t.Fatalf("Test failed. Expected: 2.72; Actual: %v", p1)
	}

	if e != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", e)
	}

	if p2 != 10.88 {
		t.Fatalf("Test failed. Expected: 10.88; Actual: %v", p2)
	}
}

func TestCalculateGradients(t *testing.T) {
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

	_, _, hist := h.CalculateGradients(img)

	target := []float32{194.86472, 280.97842, 352.12576, 246.9297, 252.13818, 51.715504, 0, 0, 0}

	if hist == nil {
		t.Fatal("Test failed")
	}

	for i := range target {
		if hist[i] != target[i] {
			t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, hist, target)
		}
	}
}
