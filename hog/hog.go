package hog

import (
	"fmt"
	"image"
)

type HOG struct {
}

func (h *HOG) CalculateGradient(img image.Image, template any) any {
	ts := len(template.([]any))

	fmt.Println(ts)

	return nil
}
