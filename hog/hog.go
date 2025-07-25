package hog

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"

	"golang.org/x/image/draw"
)

type HOG struct {
	numberOfBins int
	stepSize     int
	epsilon      float64
}

func NewHOG(numberOfBins *int, epsilon *float64) *HOG {
	instance := &HOG{
		numberOfBins: 9,
		stepSize:     20,
		epsilon:      1e-05,
	}

	if numberOfBins != nil {
		instance.numberOfBins = *numberOfBins
		instance.stepSize = 180 / instance.numberOfBins
	}

	if epsilon != nil {
		instance.epsilon = *epsilon
	}

	return instance
}

func (h *HOG) MagnitudeTheta(img [][]float32) ([][]float32, [][]float32) {
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

func (f *HOG) ImgToGray(img image.Image) *image.Gray {
	grayImg := image.NewGray(img.Bounds())

	draw.Draw(grayImg, grayImg.Bounds(), img, img.Bounds().Min, draw.Src)

	return grayImg
}

func (f *HOG) Resize(img image.Image, width, height int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.NearestNeighbor.Scale(newImg, newImg.Rect, img, img.Bounds(), draw.Over, nil)

	return newImg
}

func (f *HOG) ImgToArray(img image.Gray) [][]float32 {
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

func (f *HOG) ArrayToImg(data [][]float32, divisor *float32) (image.Image, error) {
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

func (f *HOG) CalculateJ(angle float32) float32 {
	temp := (angle / float32(f.stepSize)) - 0.5

	j := math.Floor(float64(temp))

	return float32(j)
}

func (f *HOG) CalculateCJ(j float32) float32 {
	return float32(f.stepSize) * (j + 0.5)
}

func (f *HOG) CalculateValueJ(magnitude, angle, j float32) float32 {
	Cj := f.CalculateCJ(j + 1)
	Vj := magnitude * ((Cj - angle) / float32(f.stepSize))

	return Vj
}

func (f *HOG) Partition(data [][]float32, y, x, step int) [][]float32 {
	result := make([][]float32, step)

	for i := range step {
		result[i] = make([]float32, step)
		for j := range step {
			result[i][j] = data[y+i][x+j]
		}
	}

	return result
}

func (f *HOG) BuildRow(magnitude, angle float32) (int, float32, float32) {
	valueJ := f.CalculateJ(angle)
	Vj := f.CalculateValueJ(magnitude, angle, valueJ)
	Vj_1 := magnitude - Vj

	return int(valueJ), Vj, Vj_1
}

func (f *HOG) BuildBin(magnitudes, angles [][]float32, i, j, step int) []float32 {
	var bin []float32

	magnitudeValues := f.Partition(magnitudes, i, j, step)
	angleValues := f.Partition(angles, i, j, step)

	for k := range len(magnitudeValues) {
		for l := range len(magnitudeValues[0]) {
			bin = make([]float32, f.numberOfBins)

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

func (f *HOG) HistogramPointsNine(magnitudes, angles [][]float32) [][][]float32 {
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

func (f *HOG) FetchHistValues(hist [][][]float32, i, j int) [][][]float32 {
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

func (f *HOG) CalculateK(finalVector []float32) float32 {
	var k float64

	for _, x := range finalVector {
		k += math.Pow(float64(x), 2)
	}

	k = math.Sqrt(k)

	return float32(k)
}

func (f *HOG) CalculateV2(finalVector []float32, k float32) []float32 {
	result := make([]float32, len(finalVector))

	for i, x := range finalVector {
		result[i] = x / (k + float32(f.epsilon))
	}

	return result
}

func (f *HOG) CreateFeatures(hist [][][]float32) [][][]float32 {
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

func (h *HOG) FlattenArray(data any) []float32 {
	result := make([]float32, 0)

	switch v := data.(type) {
	case [][][][]float32:
		for _, elem := range v {
			result = append(result, h.FlattenArray(elem)...)
		}
	case [][][]float32:
		for _, elem := range v {
			result = append(result, h.FlattenArray(elem)...)
		}
	case [][]float32:
		for _, elem := range v {
			result = append(result, h.FlattenArray(elem)...)
		}
	case []float32:
		for _, elem := range v {
			result = append(result, h.FlattenArray(elem)...)
		}
	case float32:
		result = append(result, v)
	}

	return result
}

func (f *HOG) HOG(img image.Image, debug bool) (image.Image, []float32, error) {
	var hogImg image.Image
	var features []float32

	resizedImg := f.Resize(img, 64, 128)

	if debug {
		filename := "outputResized.jpg"

		outputFile, err := os.Create(filename)
		if err != nil {
			return nil, nil, err
		}
		defer func() {
			outputFile.Close()
		}()

		err = jpeg.Encode(outputFile, resizedImg, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	grayImg := f.ImgToGray(resizedImg)

	if debug {
		filename := "outputGray.jpg"

		outputFile, err := os.Create(filename)
		if err != nil {
			return nil, nil, err
		}
		defer func() {
			outputFile.Close()
		}()

		err = jpeg.Encode(outputFile, grayImg, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	dump := f.ImgToArray(*grayImg)

	if debug {
		payload, _ := json.MarshalIndent(dump, "", "  ")

		os.WriteFile("outputDump.json", payload, 0644)
	}

	magnitudes, angles := f.MagnitudeTheta(dump)

	if debug {
		payload, _ := json.MarshalIndent(magnitudes, "", "  ")

		os.WriteFile("outputMagnitudes.json", payload, 0644)

		payload, _ = json.MarshalIndent(angles, "", "  ")

		os.WriteFile("outputAngles.json", payload, 0644)
	}

	histogram := f.HistogramPointsNine(magnitudes, angles)

	if debug {
		payload, _ := json.MarshalIndent(histogram, "", "  ")

		os.WriteFile("outputHist.json", payload, 0644)
	}

	features = f.FlattenArray(f.CreateFeatures(histogram))

	if debug {
		payload, _ := json.MarshalIndent(features, "", "  ")

		os.WriteFile("outputFeatures.json", payload, 0644)
	}

	hogImg, err := f.ArrayToImg(magnitudes, nil)
	if err != nil {
		return nil, nil, err
	}

	if debug {
		filename := "outputHOG.jpg"

		outputFile, err := os.Create(filename)
		if err != nil {
			return nil, nil, err
		}
		defer func() {
			outputFile.Close()
		}()

		err = jpeg.Encode(outputFile, hogImg, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	return hogImg, features, nil
}
