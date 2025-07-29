package utils_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/hog/utils"
)

func TestDistribution(t *testing.T) {
	target := `Distribution
==========================
id                   44446
gender                   5
masterCategory           7
subCategory             45
articleType            143
baseColour              47
season                   6
year                    15
usage                   30
productDisplayName   31115
`

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "styles.csv"))
	if err != nil {
		t.Fatal(err)
	}

	_, result := utils.Distribution(content)

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}

func TestHighestTen(t *testing.T) {
	classes := map[string]map[string]int{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "classes.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &classes)
	if err != nil {
		t.Fatal(err)
	}

	result := utils.HighestTen(classes)

	target := map[string]map[string]int{}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "selects.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	for key, row := range target {
		for label, count := range row {
			if count > 1 {
				if result[key][label] != target[key][label] {
					t.Fatalf(`Test failed. 
Exptected: %#v; 
Actual: %#v`, target[key], result[key])
				}
			}
		}
	}
}

func TestPlotBarGraph(t *testing.T) {
	data := map[string]map[string]int{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "selects.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	for title := range data {
		defer func() {
			os.Remove(fmt.Sprintf("%s.png", title))
		}()

		err = utils.PlotBarGraph(title, data[title], fmt.Sprintf("%s.png", title))
		if err != nil {
			t.Fatal(err)
		}
	}
}
