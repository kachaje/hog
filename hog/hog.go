package hog

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
	"sync"
)

// Img
type Img struct {
	Image image.Image
	Name  string
}

// ImageInfo image information.
type ImageInfo struct {
	Wg sync.WaitGroup
	sync.RWMutex
	Format    string
	Name      string
	Bounds    image.Rectangle
	ScaleSize int
	CellSize  int
}

// Constant
const (
	FULLCIRCLE float64 = 360
	HALFCIRCLE float64 = 180
	K          float64 = 8
	PI         float64 = math.Pi
)

// NewImageInfo return ImageInfo struct.
func NewImageInfo(format, name string, bounds image.Rectangle, scaleSize, cellSize int) *ImageInfo {
	return &ImageInfo{
		Wg:        sync.WaitGroup{},
		Format:    format,
		Name:      name,
		Bounds:    bounds,
		ScaleSize: scaleSize,
		CellSize:  cellSize,
	}
}

// Grayscale gray scale image
func (i *ImageInfo) Grayscale(imgsrc image.Image) image.Image {
	if imgsrc.ColorModel() == color.GrayModel {
		return imgsrc
	}
	bounds := imgsrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(bounds)
	for x := range w {
		for y := range h {
			gray.Set(x, y, color.GrayModel.Convert(imgsrc.At(x, y)))
		}
	}

	return gray
}

// Save save image into directory
func (i *ImageInfo) Save(filename string, imgsrc image.Image) error {
	parts := strings.Split(filename, ".")

	if len(parts) < 2 {
		return fmt.Errorf("invalid filename %s", filename)
	}

	format := parts[1]

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	switch format {
	case "png":
		return png.Encode(out, imgsrc)
	case "jpg", "jpeg":
		return jpeg.Encode(out, imgsrc, nil)
	}

	return nil
}

// Magnitude calculate the magnitude of two points
//
//	f(x, y) = sqrt(c^2 + y^2)
func Magnitude(x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

// OrientationXY calculate orientation of two points
//
//	f(x, y) = atan2(x, y) * 180 / 3.14 % 360
func OrientationXY(x, y float64) float64 {
	return math.Mod((math.Atan2(x, y) * HALFCIRCLE / math.Pi), FULLCIRCLE)
}

func xFromAngle(x, length int, angle float64) float64 {
	return math.Round(float64(x) + (float64(length) * math.Cos(angle)))
}

func yFromAngle(y, length int, angle float64) float64 {
	return math.Round(float64(y) + (float64(length) * math.Sin(angle)))
}

// Divide split rectangle into s*s cell.
func Divide(bounds image.Rectangle, s int) []image.Rectangle {
	w, h, c := bounds.Max.X, bounds.Max.Y, 0
	cells := make([]image.Rectangle, int(w/s*h/s))
	for y := 16; y < h; y += s {
		for x := 16; x < w; x += s {
			v, z := x, y
			cells[c] = image.Rect(v-s, z-s, x, y)
			c++
		}
	}
	return cells
}

// DrawLine draw a line in image.
func DrawLine(p image.Point, angle float64, length int, imgsrc image.Image, c color.Color) *image.RGBA {
	bound := imgsrc.Bounds()
	dstimg, mask := image.NewRGBA(bound), image.NewRGBA(bound)
	x1 := xFromAngle(p.X, int(length), angle)
	y1 := yFromAngle(p.Y, int(length), angle)
	x2 := xFromAngle(p.X, int(length), angle+180)
	a := (x1 - float64(p.X)) / (y1 - float64(p.Y))
	b := int(float64(p.Y) - a*float64(p.X))
	s, e := x2, x1
	if x1 < 0 {
		s, e = x1, x2
	}
	for x := int(s); x <= int(e); x++ {
		mask.Set(x, int(a*float64(x))+b, c)
	}
	draw.Draw(dstimg, bound, imgsrc, bound.Min, draw.Src)
	draw.Draw(dstimg, bound, mask, bound.Min, draw.Over)
	return dstimg
}

// HogVect hog implementation.
func HogVect(imgsrc image.Image, img *ImageInfo, debug bool) image.Image {
	bound := imgsrc.Bounds()
	hogimg := image.NewRGBA(bound)

	draw.Draw(hogimg, bound, &image.Uniform{color.Black}, image.Pt(0, 0), draw.Src)

	cells := Divide(bound, img.CellSize)

	midcell := image.Pt(int(img.CellSize/2)+1, int(img.CellSize/2)+1)
	vect := int(img.CellSize / 2)

	c := color.White //color.RGBA{0xff, 0xff, 0xff, 0xff}

	if debug {
		fmt.Printf("+ There are %d cells\n", len(cells)-1)
	}

	for k, cell := range cells {
		if cells[k] == image.Rect(0, 0, 0, 0) {
			if debug {
				log.Printf("Cell out of bound with: %d cell(s)\n", len(cells)-k)
			}
			break
		}
		img.Wg.Add(1)

		if debug {
			fmt.Printf("- Processing with %d cell\r", k)
		}

		imgcell := image.NewRGBA(cell)
		for y := cell.Min.Y; y < cell.Max.Y; y++ {
			for x := cell.Min.X; x < cell.Max.X; x++ {
				yd := math.Abs(float64(imgsrc.At(x, y-1).(color.Gray).Y - imgsrc.At(x, y+1).(color.Gray).Y))
				xd := math.Abs(float64(imgsrc.At(x-1, y).(color.Gray).Y - imgsrc.At(x+1, y).(color.Gray).Y))
				magnitude, orientation := Magnitude(xd, yd), OrientationXY(xd, yd)
				if int(magnitude)%16 == 0 { // useful i supose so!
					imgcell = DrawLine(cell.Sub(midcell).Max, orientation, vect, imgcell, c)
				}
			}

		}

		draw.Draw(hogimg, imgcell.Bounds(), imgcell, cell.Min, draw.Over)
		img.Wg.Done()
	}

	if debug {
		fmt.Printf("\n")
	}

	return hogimg
}
