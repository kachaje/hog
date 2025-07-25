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
	epsilon      = 1e-05
)

func (f *Features) ImgToArray(img image.Gray) [][]float32 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixelArray := make([][]float32, height)

	for y := range height {
		pixelArray[y] = make([]float32, width)

		for x := range width {
			pixelArray[y][x] = float32(img.At(x, y).(color.Gray).Y) / 257.0
		}
	}

	return pixelArray
}

func (f *Features) ArrayToImg(data [][]float32, divisor *float32) (image.Image, error) {
	factor := float32(257.0)
	if divisor != nil {
		factor = *divisor
	}

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
			value := uint8(data[y][x] * factor)

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

func (f *Features) CalculateCJ(j float32) float32 {
	return float32(stepSize) * (j + 0.5)
}

func (f *Features) CalculateValueJ(magnitude, angle, j float32) float32 {
	Cj := f.CalculateCJ(j + 1)
	Vj := magnitude * ((Cj - angle) / float32(stepSize))

	return Vj
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

func (f *Features) BuildRow(magnitude, angle float32) (int, float32, float32) {
	valueJ := f.CalculateJ(angle)
	Vj := f.CalculateValueJ(magnitude, angle, valueJ)
	Vj_1 := magnitude - Vj

	return int(valueJ), Vj, Vj_1
}

func (f *Features) BuildBin(magnitudes, angles [][]float32, i, j, step int) []float32 {
	var bin []float32

	magnitudeValues := f.Partition(magnitudes, i, j, step)
	angleValues := f.Partition(angles, i, j, step)

	for k := range len(magnitudeValues) {
		for l := range len(magnitudeValues[0]) {
			bin = make([]float32, numberOfBins)

			valueJ, Vj, Vj_1 := f.BuildRow(magnitudeValues[k][l], angleValues[k][l])

			if valueJ < 0 {
				bin[step] += Vj
				bin[0] += Vj_1
			} else {
				bin[valueJ] += Vj
				bin[valueJ+1] += Vj_1
			}
		}
	}

	return bin
}

func (f *Features) HistogramPointsNine(magnitudes, angles [][]float32) [][][]float32 {
	hist := make([][][]float32, 0)

	step := 8
	height := len(magnitudes)
	width := len(magnitudes[0])

	for i := 0; i < height; i += step {
		temp := make([][]float32, 0)

		for j := 0; j < width; j += step {
			bins := f.BuildBin(magnitudes, angles, i, j, step)

			temp = append(temp, bins)
		}

		hist = append(hist, temp)
	}

	return hist
}

func (f *Features) FetchHistValues(hist [][][]float32, i, j int) [][][]float32 {
	values := make([][][]float32, 0)

	for k := range 2 {
		row := [][]float32{}

		for l := range 2 {
			row = append(row, hist[i+k][j+l])
		}

		values = append(values, row)
	}

	return values
}

func (f *Features) CalculateK(finalVector []float32) float32 {
	var k float64

	for _, x := range finalVector {
		k += math.Pow(float64(x), 2)
	}

	k = math.Sqrt(k)

	return float32(k)
}

func (f *Features) CalculateV2(finalVector []float32, k float32) []float32 {
	result := make([]float32, len(finalVector))

	for i, x := range finalVector {
		result[i] = x / (k + float32(epsilon))
	}

	return result
}

func (f *Features) CreateFeatures(hist [][][]float32) [][][]float32 {
	featureVectors := [][][]float32{}
	epsilon := 1e-05

	_ = epsilon

	for i := range len(hist) - 1 {
		temp := [][]float32{}
		for j := range len(hist[0]) - 1 {
			values := f.FetchHistValues(hist, i, j)

			finalVector := []float32{}
			for _, k := range values {
				for _, l := range k {
					finalVector = append(finalVector, l...)
				}
			}

			k := f.CalculateK(finalVector)

			finalVector = f.CalculateV2(finalVector, k)

			temp = append(temp, finalVector)
		}
		featureVectors = append(featureVectors, temp)
	}

	return featureVectors
}
