package utils

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
)

func CleanString(str string) string {
	return regexp.MustCompile(`\t|\n|\s`).ReplaceAllLiteralString(str, "")
}

func Distribution(content []byte) (map[string][]string, string) {
	rows := func() [][]string {
		result := make([][]string, 0)

		for rowData := range strings.SplitSeq(string(content), "\n") {
			row := []string{}

			for text := range strings.SplitSeq(rowData, ",") {
				if text != "" {
					row = append(row, text)
				}
			}

			if len(row) > 0 {
				result = append(result, row)
			}
		}

		return result
	}()

	selectClasses := []string{"masterCategory", "subCategory", "season", "gender", "usage"}

	labels := map[int]string{}
	counts := map[string][]string{}
	classes := map[string]map[string]int{}

	for i, row := range rows {
		if i == 0 {
			for j, text := range row {
				if text == "" {
					continue
				}

				labels[j] = text
			}
		} else {
			for j, text := range row {
				if text == "" {
					continue
				}

				label := labels[j]

				if slices.Contains(selectClasses, label) {
					if classes[label] == nil {
						classes[label] = map[string]int{}
					}

					classes[label][text]++
				}

				if counts[label] == nil || !slices.Contains(counts[label], text) {
					counts[label] = append(counts[label], text)
				}
			}
		}
	}

	stats := []string{
		"Distribution",
		"==========================",
	}

	for i := range len(labels) {
		key := labels[i]

		if key == "" {
			continue
		}

		value := counts[key]

		stats = append(stats, fmt.Sprintf("%-20s%6d", key, len(value)))
	}

	return counts, strings.Join(stats, "\n")

}

func HighestTen(data map[string]map[string]int) map[string]map[string]int {
	result := map[string]map[string]int{}

	for key, row := range data {
		result[key] = map[string]int{}

		rows := map[int]string{}

		for key, value := range row {
			rows[value] = key
		}

		indices := []int{}

		for i := range rows {
			indices = append(indices, i)
		}

		sort.Sort(sort.Reverse(sort.IntSlice(indices)))

		limit := min(len(indices), 10)

		for _, index := range indices[:limit] {
			result[key][rows[index]] = index
		}
	}

	return result
}

func PlotBarGraph(title string, data map[string]int, filename string) error {
	labels := []string{}

	type kv struct {
		Key   string
		Value int
	}
	var ss []kv

	for k, v := range data {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	items := make([]opts.BarData, 0)
	for _, val := range ss {
		labels = append(labels, val.Key)
		items = append(items, opts.BarData{Value: val.Value})
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			BackgroundColor: "#FFFFFF",
		}),
		// Don't forget disable the Animation
		charts.WithAnimation(false),
		charts.WithTitleOpts(opts.Title{
			Title:     title,
			TextAlign: "center",
			Left:      "50%",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Rotate: 30,
			},
		}),
	)
	bar.SetXAxis(labels).
		AddSeries("", items)

	render.MakeChartSnapshot(bar.RenderContent(), filename)

	return nil
}
