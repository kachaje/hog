package hog

import (
	"image"
	"image/draw"
)

type HOG struct{}

func (h *HOG) ImgToGray(img image.Image) *image.Gray {
	grayImg := image.NewGray(img.Bounds())

	draw.Draw(grayImg, grayImg.Bounds(), img, img.Bounds().Min, draw.Src)

	return grayImg
}
