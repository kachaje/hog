package hog_test

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"image/color"
	"image/draw"
	_ "image/jpeg"

	"github.com/kachaje/hog/hog"
)

func TestSumMatrix(t *testing.T) {
	h := hog.NewHOG(nil)

	arr := []any{
		1.0,
		2.0,
		[]any{3.0, 4.0, []any{5.0, 6.0}},
		7.0,
		[]any{8.0, []any{9.0}},
	}

	result := h.SumMatrix(arr)

	if result != 45 {
		t.Fatalf("Test failed. Expected: 45; Actual: %v", result)
	}
}

func TestMultiplyMatrices(t *testing.T) {
	h := hog.NewHOG(nil)

	A := [][]float32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	B := [][]float32{{9, 8}, {6, 5}, {3, 2}}

	result, err := h.MultiplyMatrices(A, B)
	if err != nil {
		t.Fatal(err)
	}

	target := [][]float32{{30, 24}, {84, 69}, {138, 114}}

	if result == nil {
		t.Fatal("Test failed")
	}

	for i := range target {
		for j := range target[i] {
			if result[i][j] != target[i][j] {
				t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, result, target)
			}
		}
	}
}

func TestGetRegion(t *testing.T) {
	width := 3
	height := 3

	grayImg := image.NewGray(image.Rect(0, 0, width, height))

	i := 0
	for y := range height {
		for x := range width {
			i++
			value := uint8(i * 255)
			grayImg.SetGray(x, y, color.Gray{value})
		}
	}

	h := hog.NewHOG(grayImg)

	result := h.GetRegion(1, 1, 3)

	target := [][]float32{
		{
			0.9922179, 0.98054475, 0.9688716,
		},
		{
			0.98832685, 0.9766537, 0.96498054,
		},
		{
			0.9844358, 0.97276264, 0.9610895,
		},
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	for i := range target {
		for j := range target[i] {
			if result[i][j] != target[i][j] {
				t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, result, target)
			}
		}
	}
}

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

	grayImg := image.NewGray(img.Bounds())

	draw.Draw(grayImg, grayImg.Bounds(), img, img.Bounds().Min, draw.Src)

	h := hog.NewHOG(grayImg)

	template := [][]float32{{-1}, {0}, {1}}

	result := h.CalculateGradient(template)

	if len(result) != 429 {
		t.Fatalf("Test failed. Expected: 429; Actual: %v", len(result))
	}
}
