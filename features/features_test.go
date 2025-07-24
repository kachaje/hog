package features_test

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
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

	data, err := os.ReadFile("./fixtures/dump.json")
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
	t.Skip()

	var targetData, magData, thetaData [][]float32

	data, err := os.ReadFile("./fixtures/dump.json")
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

		if false {
			os.Remove(magFilename)
		}
	}()

	magImg, err := f.ArrayToImg(magData, nil)
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

		if false {
			os.Remove(thetaFilename)
		}
	}()

	factor := float32(math.Pi * 257.0 / 180)
	thetaImg, err := f.ArrayToImg(thetaData, &factor)
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

func TestCalculateCJ(t *testing.T) {
	f := features.Features{}

	targets := map[float32]float32{
		4: 90.0,
		3: 70.0,
		2: 50.0,
		1: 30.0,
		0: 10.0,
	}

	for j, value := range targets {
		result := f.CalculateCJ(j)

		if result != value {
			t.Fatalf("Test failed. Expected: %v; Actual: %v\n", value, result)
		}
	}
}

func TestCalculateValueJ(t *testing.T) {
	f := features.Features{}

	targets := map[float32]map[string]float32{
		0.002121697: {
			"angle": 40.989482415, "Vj": 0.00095587934, "valueJ": 1,
		},
		0.005973429: {
			"angle": 51.403386204, "Vj": 0.005554278, "valueJ": 2,
		},
		0.034484309: {
			"angle": 24.476863127, "Vj": 0.009523077, "valueJ": 0,
		},
		0.029104762: {
			"angle": 88.216660375, "Vj": 0.0025951848, "valueJ": 3,
		},
		0.029155932: {
			"angle": 88.555655828, "Vj": 0.0021055592, "valueJ": 3,
		},
		0.0305682: {
			"angle": 83.918670745, "Vj": 0.009294765, "valueJ": 3,
		},
		0.032742355: {
			"angle": 83.564826236, "Vj": 0.010535136, "valueJ": 3,
		},
		0.03223796: {
			"angle": 88.745362256, "Vj": 0.0020223495, "valueJ": 3,
		},
		0.03359077: {
			"angle": 69.121575212, "Vj": 0.0014753497, "valueJ": 2,
		},
		0.049294209: {
			"angle": 46.838597471, "Vj": 0.007791944, "valueJ": 1,
		},
	}

	for magnitude, row := range targets {
		angle := row["angle"]
		vj := row["Vj"]
		valueJ := row["valueJ"]

		result := f.CalculateValueJ(magnitude, angle, valueJ)

		if result != vj {
			t.Fatalf("Test failed on %v. Expected: %v; Actual: %v\n", magnitude, vj, result)
		}
	}
}

func TestBuildRow(t *testing.T) {
	f := features.Features{}

	files, err := os.ReadDir("./fixtures/data/points")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(".", "fixtures", "data", "points", file.Name()))
		if err != nil {
			t.Fatal(err)
		}

		var data map[string]float32

		err = json.Unmarshal(content, &data)
		if err != nil {
			t.Fatal(err)
		}

		magnitude, angle, targetVj, targetVj_1, targetValueJ := data["magnitude"], data["angle"], data["Vj"], data["Vj_1"], data["value_j"]

		valueJ, Vj, Vj_1 := f.BuildRow(magnitude, angle)

		Vj = float32(math.Floor(float64(Vj*1e6)) / 1e6)
		targetVj = float32(math.Floor(float64(targetVj*1e6)) / 1e6)

		Vj_1 = float32(math.Floor(float64(Vj_1*1e6)) / 1e6)
		targetVj_1 = float32(math.Floor(float64(targetVj_1*1e6)) / 1e6)

		if Vj != targetVj {
			t.Fatalf("Test failed. Expected: %v; Actual: %v", targetVj, Vj)
		}
		if Vj_1 != float32(targetVj_1) {
			t.Fatalf("Test failed. Expected: %v; Actual: %v", targetVj_1, Vj_1)
		}

		if int(targetValueJ) != valueJ {
			t.Fatalf("Test failed. Expected: %v; Actual: %v", targetValueJ, valueJ)
		}
	}
}

func TestBuildBin(t *testing.T) {
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

	bin := f.BuildBin(magData, thetaData, 0, 0, 8)

	fmt.Println(bin)
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

	if len(hist) != 16 {
		t.Fatalf("Test failed. Expected: 16; Actual: %v", len(hist))
	}

	if len(hist[0]) != 8 {
		t.Fatalf("Test failed. Expected: 8; Actual: %v", len(hist[0]))
	}

	if len(hist[0][0]) != 9 {
		t.Fatalf("Test failed. Expected: 9; Actual: %v", len(hist[0][0]))
	}

	if false {
		payload, _ := json.Marshal(hist)

		os.WriteFile("./fixtures/hist.json", payload, 0644)
	}
}
