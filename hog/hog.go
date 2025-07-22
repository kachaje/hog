package hog

import (
	"fmt"
	"image"
)

type HOG struct {
}

func (h *HOG) CalculateGradient(img image.Image, template any) any {
	ts := len(template.([]any))

	bounds := img.Bounds()

	rect := image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(bounds.Max.X+ts-1, bounds.Max.Y+ts-1),
	}

	newImg := image.NewRGBA(rect)

	fmt.Println(ts, newImg.Rect)

	return nil
}
