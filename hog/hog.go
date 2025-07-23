package hog

import (
	"image"
	"image/color"
	"image/draw"
)

type HOG struct{}

func (h *HOG) ImgToGray(img image.Image) *image.Gray {
	grayImg := image.NewGray(img.Bounds())

	draw.Draw(grayImg, grayImg.Bounds(), img, img.Bounds().Min, draw.Src)

	return grayImg
}

func (h *HOG) ResizeImg(img image.Image, height int) image.Image {
	if height < 50 {
		return img
	}
	bounds := img.Bounds()
	imgHeight := bounds.Dy()
	if height >= imgHeight {
		return img
	}
	imgWidth := bounds.Dx()
	resizeFactor := float32(imgHeight) / float32(height)
	ratio := float32(imgWidth) / float32(imgHeight)
	width := int(float32(height) * ratio)
	resizedImage := image.NewRGBA(image.Rect(0, 0, width, height))

	var imgX, imgY int
	var imgColor color.Color
	for x := range width {
		for y := range height {
			imgX = int(resizeFactor*float32(x) + 0.5)
			imgY = int(resizeFactor*float32(y) + 0.5)
			imgColor = img.At(imgX, imgY)
			resizedImage.Set(x, y, imgColor)
		}
	}

	return resizedImage
}
