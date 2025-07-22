package hog

import (
	"fmt"
	"image"
	"image/color"
)

type HOG struct {
	grayImg image.Image
}

func NewHOG(grayImg image.Image) *HOG {
	return &HOG{
		grayImg: grayImg,
	}
}

func (h *HOG) Multiply(A, B [][]float32) ([][]float32, error) {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	if colsA != rowsB {
		return nil, fmt.Errorf("matrix multiplication not valid")
	}

	C := make([][]float32, rowsA)
	for i := range C {
		C[i] = make([]float32, colsB)
	}

	for i := range rowsA {
		for j := range colsB {
			for k := range colsA {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}

	return C, nil
}

func (h *HOG) GetRegion(r, c, step int) [][]float32 {
	result := make([][]float32, 3)

	for i := range step {
		result[i] = make([]float32, 3)

		for j := range step {
			value := float32(h.grayImg.At(r+i-1, c+j-1).(color.Gray).Y) / 257.0
			result[i][j] = value
		}
	}

	return result
}

func (h *HOG) CalculateGradient(template any) image.Image {
	step := len(template.([]any))

	bounds := h.grayImg.Bounds()

	rect := image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(bounds.Max.X+step-1, bounds.Max.Y+step-1),
	}

	newImg := image.NewRGBA(rect)
	result := image.NewRGBA(rect)

	_ = newImg

	for r := range h.grayImg.Bounds().Max.X {
		for c := range h.grayImg.Bounds().Max.Y {
			currentRegion := h.GetRegion(r, c, step)

			fmt.Printf("(%v, %v) - %v\n", r, c, currentRegion)
		}
	}

	return result
}
