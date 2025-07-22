package hog

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"reflect"
)

type HOG struct {
	grayImg image.Image
}

func NewHOG(grayImg image.Image) *HOG {
	return &HOG{
		grayImg: grayImg,
	}
}

func (h *HOG) SumMatrix(arr any) float32 {
	var result float32

	if reflect.TypeOf(arr).Kind() == reflect.Slice {
		s := reflect.ValueOf(arr)

		for i := range s.Len() {
			element := s.Index(i).Interface()

			if val, ok := element.(float32); ok {
				result += val
			} else if val, ok := arr.(float64); ok {
				result += float32(val)
			} else {
				result += h.SumMatrix(element)
			}
		}
	} else if val, ok := arr.(float32); ok {
		result += val
	} else if val, ok := arr.(float64); ok {
		result += float32(val)
	}

	return result
}

func (h *HOG) MultiplyMatrices(A, B [][]float32) ([][]float32, error) {
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

func (h *HOG) CalculateGradient(template [][]float32) [][]float32 {
	step := len(template)

	bounds := h.grayImg.Bounds()

	rect := image.Rectangle{
		Max: image.Pt(bounds.Max.X+step-1, bounds.Max.Y+step-1),
	}

	result := make([][]float32, rect.Max.Y)
	for i := range rect.Max.Y {
		result[i] = make([]float32, rect.Max.X)
	}

	for r := range h.grayImg.Bounds().Max.Y {
		for c := range h.grayImg.Bounds().Max.X {
			currentRegion := h.GetRegion(r, c, step)

			currentResult, err := h.MultiplyMatrices(currentRegion, template)
			if err != nil {
				log.Println(err)
				continue
			}

			score := h.SumMatrix(currentResult)

			result[r][c] = score
		}
	}

	return result
}
