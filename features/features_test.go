package features_test

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"testing"

	_ "image/jpeg"
	"image/png"

	"github.com/kachaje/hog/features"
	"github.com/kachaje/hog/hog"
)

func TestImgToArray(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "flowerGray.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	h := hog.HOG{}
	f := features.Features{}

	grayImg := h.ImgToGray(img)

	result := f.ImgToArray(*grayImg)

	if result == nil {
		t.Fatal("Test failed")
	}

	var target [][]float32

	data, err := os.ReadFile("../data/dump.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &target)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	for i := range target {
		for j := range target[i] {
			if result[i][j] != target[i][j] {
				t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, result, target)
			}
		}
	}
}

func TestMagnitudeTheta(t *testing.T) {
	var targetData, magData, thetaData [][]float32

	data, err := os.ReadFile("../data/dump.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile("./fixtures/magnitudes.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &magData)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile("./fixtures/thetas.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &thetaData)
	if err != nil {
		t.Fatal(err)
	}

	f := features.Features{}

	mag, theta := f.MagnitudeTheta(targetData)

	if mag == nil {
		t.Fatal("Test failed")
	}

	for i := range magData {
		for j := range magData[i] {
			if mag[i][j] != magData[i][j] {
				t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, magData, mag)
			}
		}
	}

	if theta == nil {
		t.Fatal("Test failed")
	}

	for i := range thetaData {
		for j := range thetaData[i] {
			if theta[i][j] != thetaData[i][j] {
				t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, thetaData, theta)
			}
		}
	}
}

func TestArrayToImg(t *testing.T) {
	var magData, thetaData [][]float32

	data, err := os.ReadFile("./fixtures/magnitudes.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &magData)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile("./fixtures/thetas.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &thetaData)
	if err != nil {
		t.Fatal(err)
	}

	f := features.Features{}

	magFilename := "outputMag.png"
	thetaFilename := "outputTheta.png"

	magFile, err := os.Create(magFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		magFile.Close()

		os.Remove(magFilename)
	}()

	magImg, err := f.ArrayToImg(magData)
	if err != nil {
		t.Fatal(err)
	}

	if err := png.Encode(magFile, magImg); err != nil {
		t.Fatal(err)
	}

	thetaFile, err := os.Create(thetaFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		magFile.Close()

		os.Remove(thetaFilename)
	}()

	thetaImg, err := f.ArrayToImg(thetaData)
	if err != nil {
		t.Fatal(err)
	}

	if err := png.Encode(thetaFile, thetaImg); err != nil {
		t.Fatal(err)
	}
}

func TestPartition(t *testing.T) {
	data := [][]float32{
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{10, 11, 12, 13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24, 25, 26, 27},
		{28, 29, 30, 31, 32, 33, 34, 35, 36},
		{37, 38, 39, 40, 42, 43, 44, 45, 46},
		{47, 48, 49, 50, 51, 52, 53, 54, 55},
	}
	target := [][]float32{
		{12, 13, 14},
		{21, 22, 23},
		{30, 31, 32},
	}

	step := 3

	f := features.Features{}

	result := f.Partition(data, 1, 2, step)

	if result == nil {
		t.Fatal("Test failed")
	}

	for i := range len(target) {
		for j := range len(target[0]) {
			if result[i][j] != target[i][j] {
				t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, target, result)
			}
		}
	}

	targets := [][][]float32{
		{
			{1, 2, 3},
			{10, 11, 12},
			{19, 20, 21},
		},
		{
			{4, 5, 6},
			{13, 14, 15},
			{22, 23, 24},
		},
		{
			{7, 8, 9},
			{16, 17, 18},
			{25, 26, 27},
		},
		{
			{28, 29, 30},
			{37, 38, 39},
			{47, 48, 49},
		},
		{
			{31, 32, 33},
			{40, 42, 43},
			{50, 51, 52},
		},
		{
			{34, 35, 36},
			{44, 45, 46},
			{53, 54, 55},
		},
	}

	k := 0
	for i := 0; i < len(data); i += step {
		for j := 0; j < len(data[0]); j += step {
			result := f.Partition(data, i, j, step)

			if result == nil {
				t.Fatal("Test failed")
			}

			target := targets[k]

			for i := range len(target) {
				for j := range len(target[0]) {
					if result[i][j] != target[i][j] {
						t.Fatalf(`Test failed. 
Expected: %#v; 
Actual: %#v`, target, result)
					}
				}
			}

			k++
		}
	}

	_ = targets
}

func TestCalculateJ(t *testing.T) {
	f := features.Features{}

	targets := map[float32]float32{
		89.699551773: 3,
		88.745362256: 3,
		69.121575212: 2,
		46.838597471: 1,
		24.476863127: 0,
		14.04109124:  0,
		69.857194452: 2,
		40.989482415: 1,
	}

	for angle, value := range targets {
		result := f.CalculateJ(angle)

		if result != value {
			t.Fatalf("Test failed. Expected: %v; Actual: %v\n", value, result)
		}
	}
}

func TestHistogramPointsNine(t *testing.T) {
	var magData, thetaData [][]float32

	data, err := os.ReadFile("./fixtures/magnitudes.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &magData)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile("./fixtures/thetas.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &thetaData)
	if err != nil {
		t.Fatal(err)
	}

	f := features.Features{}

	hist := f.HistogramPointsNine(magData, thetaData)

	fmt.Println(hist)
}
