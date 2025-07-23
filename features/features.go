package features

import (
	"image"
	"image/color"
	"math"
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

func (f *Features) MagnitudeTheta(img [][]float32) ([][]float32, [][]float32) {
	height := 128
	width := 64

	mag := make([][]float32, height)
	theta := make([][]float32, height)

	for i := range height {
		var Gx, Gy float32

		mag[i] = make([]float32, width)
		theta[i] = make([]float32, width)

		for j := range width {
			// Condition for axis 0
			if j-1 <= 0 || j+1 >= width {
				if j-1 <= 0 {
					// Condition if first element
					Gx = img[i][j+1] - 0
				} else if j+1 >= len(img[0]) {
					Gx = 0 - img[i][j-1]
				}
				// Condition for first element
			} else {
				Gx = img[i][j+1] - img[i][j-1]
			}

			// Condition for axis 1
			if i-1 <= 0 || i+1 >= height {
				if i-1 <= 0 {
					Gy = 0 - img[i+1][j]
				} else if i+1 >= 128 {
					Gy = img[i-1][j] - 0
				}
			} else {
				Gy = img[i-1][j] - img[i+1][j]
			}

			// Calculating magnitude
			magnitude := math.Round(math.Sqrt(math.Pow(float64(Gx), 2)+math.Pow(float64(Gy), 2))*1e9) / 1e9

			mag[i] = append(mag[i], float32(magnitude))

			var angle float64

			if Gx == 0 {
				angle = 0.0
			} else {
				angle = math.Round(math.Abs(math.Atan(float64(Gy)/float64(Gx))*180/math.Pi)*1e9) / 1e9
			}

			theta[i] = append(theta[i], float32(angle))
		}
	}

	return mag, theta
}
