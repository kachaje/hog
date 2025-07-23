package main

import (
	"log"
	"os/exec"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	pts := plotter.XYs{
		{X: 1, Y: 1},
		{X: 2, Y: 3},
		{X: 3, Y: 2},
		{X: 4, Y: 4},
		{X: 5, Y: 5},
	}

	p := plot.New()

	p.Title.Text = "Simple Scatter Plot"
	p.X.Label.Text = "X-axis"
	p.Y.Label.Text = "Y-axis"

	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}
	p.Add(s)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "output.png"); err != nil {
		log.Panic(err)
	}

	cmd := exec.Command("open", "output.png")
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
}
