package hog

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
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

func (h *HOG) GradX(img image.Image, x, y int) float32 {
	g1 := float32(img.At(x, y-1).(color.Gray).Y)
	g2 := float32(img.At(x, y+1).(color.Gray).Y)

	grad := g2 - g1

	return grad
}

func (h *HOG) GradY(img image.Image, x, y int) float32 {
	g1 := float32(img.At(x-1, y).(color.Gray).Y)
	g2 := float32(img.At(x+1, y).(color.Gray).Y)

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

func (h *HOG) AngleWeight(magnitude, degrees float32) (float32, float32, float32, float32) {
	var groupStart float32
	var groupEnd float32
	var part1 float32
	var part2 float32

	binSize := float32(20)

	groupStart = float32(int(degrees/binSize)) * binSize
	groupEnd = groupStart + binSize

	part1 = ((groupEnd - degrees) / binSize) * magnitude
	part2 = ((degrees - groupStart) / binSize) * magnitude

	return groupStart / binSize, part1, groupEnd / binSize, part2
}

func (h *HOG) CalculateGradients(img image.Image) ([][]float32, image.Image, []int) {
	var hog [][]float32
	var hogImg image.Image
	hist := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}

	for r := range img.Bounds().Max.Y {
		for c := range img.Bounds().Max.X {
			gx := h.GradX(img, c, r)
			gy := h.GradY(img, c, r)
			mag, _, deg := h.GradOrien(gx, gy)

			fmt.Println(int(mag), int(deg))
		}
	}

	return hog, hogImg, hist
}
