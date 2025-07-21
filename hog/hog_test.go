package hog_test

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/hog/hog"
)

func TestDrawLine(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "face.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}
	bounds := img.Bounds()

	result := hog.DrawLine(image.Pt(bounds.Max.X/2, bounds.Max.Y/2), 0.5, 100, img, color.RGBA{R: 255})

	if result.Bounds() != bounds {
		t.Fatal("Test failed")
	}

	filename := "outputLine.png"

	outputFile, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		outputFile.Close()

		os.Remove(filename)
	}()

	err = jpeg.Encode(outputFile, result, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMagnitude(t *testing.T) {
	result := hog.Magnitude(5, 5)

	target := 7.0710678118654755

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}

func TestDivide(t *testing.T) {
	result := hog.Divide(image.Rect(0, 0, 4, 4), 2)

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %d\n", len(result))
	}
}

func TestHogGrayscale(t *testing.T) {
	reader, err := os.Open(filepath.Join("..", "data", "face.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	i := hog.ImageInfo{}

	result := i.Grayscale(img)

	filename := "outputGray.png"

	outputFile, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		outputFile.Close()

		os.Remove(filename)
	}()

	err = jpeg.Encode(outputFile, result, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHogSave(t *testing.T) {
	filename := "outputSave.png"

	reader, err := os.Open(filepath.Join("..", "data", "face.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		reader.Close()

		_, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			os.Remove(filename)
		}
	}()

	img, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	i := hog.ImageInfo{}

	err = i.Save(filename, img)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHogMagnitude(t *testing.T) {
	result := hog.Magnitude(4, 2)

	target := 4.47213595499958

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}

func TestHogOrientationXY(t *testing.T) {
	result := hog.OrientationXY(10, 8)

	target := 51.34019174590991

	if result != target {
		t.Fatalf("Test failed. Expected: %f; Actual: %f\n", target, result)
	}
}

func TestHogVect(t *testing.T) {
	filename := "face.jpg"
	outputFilename := "outputHog.png"

	reader, err := os.Open(filepath.Join("..", "data", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		reader.Close()

		_, err := os.Stat(outputFilename)
		if !os.IsNotExist(err) && false {
			os.Remove(outputFilename)
		}
	}()

	srcImg, format, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}

	imgInfo := hog.NewImageInfo(format, filename, srcImg.Bounds(), 17, 5)

	grayImg := imgInfo.Grayscale(srcImg)

	imgHog := hog.HogVect(grayImg, imgInfo, false)

	err = imgInfo.Save(outputFilename, imgHog)
	if err != nil {
		t.Fatal(err)
	}
}
