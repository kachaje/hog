package features

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Features struct{}

var (
	numberOfBins = 9
	stepSize     = 180 / numberOfBins
)

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

		mag[i] = make([]float32, 0)
		theta[i] = make([]float32, 0)

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

func (f *Features) ArrayToImg(data [][]float32) (image.Image, error) {
	height := len(data)
	if height == 0 {
		return nil, fmt.Errorf("pixel data is empty")
	}
	width := len(data[0])
	if width == 0 {
		return nil, fmt.Errorf("inner array of pixel data is empty")
	}

	img := image.NewGray(image.Rect(0, 0, width, height))

	for y := range height {
		for x := range width {
			value := uint8(data[y][x])

			c := color.Gray{value}

			img.Set(x, y, c)
		}
	}

	return img, nil
}

func (f *Features) CalculateJ(angle float32) float32 {
	temp := (angle / float32(stepSize)) - 0.5

	j := math.Floor(float64(temp))

	return float32(j)
}

func (f *Features) Partition(data [][]float32, y, x, step int) [][]float32 {
	result := make([][]float32, step)

	for i := range step {
		result[i] = make([]float32, step)
		for j := range step {
			result[i][j] = data[y+i][x+j]
		}
	}

	return result
}

func (f *Features) HistogramPointsNine(mag, theta [][]float32) [][]float32 {
	hist := make([][]float32, 0)

	step := 8
	height := len(mag)
	width := len(mag[0])

	for i := 0; i < height; i += step {
		temp := make([][]float32, 0)
		_ = temp

		for j := 0; j < width; j += step {
			for k := range step {
				magnitudeValues := f.Partition(mag, i, j, step)
				angleValues := f.Partition(theta, i, j, step)

				_ = angleValues

				for l := range len(magnitudeValues[0]) {
					bins := make([]float32, numberOfBins)

					fmt.Println(i, j, k, l, bins)
				}
			}
		}
	}

	return hist
}
