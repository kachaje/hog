package main

import (
	"encoding/json"
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"

	"github.com/kachaje/hog/hog"
)

func main() {
	var filename string
	var debug, show bool

	flag.StringVar(&filename, "f", "", "file to work with")
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.BoolVar(&show, "s", false, "visualise image")

	flag.Parse()

	if filename == "" {
		log.Fatal("Missing required filename")
	}

	h := hog.NewHOG(nil, nil)

	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	hogImg, features, err := h.HOG(img, debug)
	if err != nil {
		log.Fatal(err)
	}

	filename = "outputHOG.png"

	outputFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		outputFile.Close()
	}()

	err = png.Encode(outputFile, hogImg)
	if err != nil {
		log.Fatal(err)
	}

	payload, err := json.MarshalIndent(features, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("outputFeatures.json", payload, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if show {
		cmd := exec.Command("open", filename)
		_, err = cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
	}
}
