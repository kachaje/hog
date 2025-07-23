package hog

import (
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
)

type HOG struct {
}

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

func (h *HOG) ResizeShrink(img image.Image, width, height int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.NearestNeighbor.Scale(newImg, newImg.Rect, img, img.Bounds(), draw.Over, nil)

	return newImg
}

func (h *HOG) ImgToArray(img image.Gray) [][]float32 {
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

func (h *HOG) GradX(img image.Gray, x, y int) float32 {
	g1 := float32(img.At(x, y-1).(color.Gray).Y) / 257.0
	g2 := float32(img.At(x, y+1).(color.Gray).Y) / 257.0

	grad := g2 - g1

	return grad
}

func (h *HOG) GradY(img image.Gray, x, y int) float32 {
	g1 := float32(img.At(x-1, y).(color.Gray).Y) / 257.0
	g2 := float32(img.At(x+1, y).(color.Gray).Y) / 257.0

	grad := g2 - g1

	return grad
}

func (h *HOG) GradOrien(gx, gy float32) (float64, float64, float64) {
	gx2 := float64(gx * gx)
	gy2 := float64(gy * gy)

	magnitude := math.Sqrt(gx2 + gy2)
	orientationRad := math.Atan(float64(gy) / float64(gx))
	orientationDeg := orientationRad * 180 / math.Pi

	return magnitude, orientationRad, orientationDeg
}

func (h *HOG) Gradient(img image.Image) image.Image {
	var newImg image.Image

	return newImg
}
