package features

import (
	"image"
	"image/color"
)

type Features struct{}

func (f *Features) ImgToArray(img image.Gray) [][]float32 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixelArray := make([][]float32, height)

	for y := range height {
		pixelArray[y] = make([]float32, width)

		for x := range width {
			pixelArray[y][x] = float32(img.At(x, y).(color.Gray).Y)
		}
	}

	return pixelArray
}
