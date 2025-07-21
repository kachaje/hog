package hog

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"sync"

	"golang.org/x/image/draw"
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
	Format   string
	Name     string
	Bounds   image.Rectangle
	Scalsize int
	Cellsize int
}

// NewImageInfo return ImageInfo struct.
func NewImageInfo(f, n string, b image.Rectangle, s, c int) *ImageInfo {
	return &ImageInfo{
		Wg:       sync.WaitGroup{},
		Format:   f,
		Name:     n,
		Bounds:   b,
		Scalsize: s,
		Cellsize: c,
	}
}

// Grayscale gray scale image
func (i *ImageInfo) Grayscale(imgsrc image.Image) image.Image {
	log.Println("+ Grascaling image ...")
	if imgsrc.ColorModel() == color.GrayModel {
		return imgsrc
	}
	bounds := imgsrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(bounds)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			gray.Set(x, y, color.GrayModel.Convert(imgsrc.At(x, y)))
		}
	}
	return gray
}

// Save save image into directory
func (i *ImageInfo) Save(name string, imgsrc image.Image) {
	out, err := os.Create(fmt.Sprintf("data/%s-%s.gore.%s", name, i.Name, i.Format))
	if err != nil {
		log.Fatalf("image.go:makeItGray:os.Create: image: %s %v", name, err)
	}
	defer out.Close()
	log.Printf("Saving %s-%s.gore.%s", name, i.Name, i.Format)
	switch i.Format {
	case "png":
		png.Encode(out, imgsrc)
	case "jpg", "jpeg":
		jpeg.Encode(out, imgsrc, nil)
	}
}

// Scale reduce image into i.scalsize defind in ImageInfo.
func (i *ImageInfo) Scale(imgsrc image.Image) image.Image {
	log.Printf("+ Scale image into %d", i.Scalsize)
	bound := imgsrc.Bounds()
	rect := image.Rect(0, 0, int(bound.Max.X/i.Scalsize), int(bound.Max.Y/i.Scalsize))
	dstimg := image.NewRGBA(rect)
	draw.ApproxBiLinear.Scale(dstimg, rect, imgsrc, imgsrc.Bounds(), draw.Over, nil)
	return dstimg
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

func RGBChannel(imgsrc image.Image, channel string) image.Image {
	if imgsrc.ColorModel() == color.GrayModel {
		return setColorGray(imgsrc, channel)
	} else {
		return setColorRGBA(imgsrc, channel)
	}
}

func setColorRGBA(imgsrc image.Image, channel string) image.Image {
	maxX, maxY := imgsrc.Bounds().Max.X, imgsrc.Bounds().Max.Y
	imgdst := image.NewRGBA(imgsrc.Bounds())
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			r, g, b := uint32(0), uint32(0), uint32(0)
			switch channel {
			case "red":
				r, _, _, _ = imgsrc.At(x, y).RGBA()
			case "green":
				_, g, _, _ = imgsrc.At(x, y).RGBA()
			case "blue":
				_, _, b, _ = imgsrc.At(x, y).RGBA()
			default:
				r, g, b, _ = imgsrc.At(x, y).RGBA()
			}
			imgdst.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return imgdst
}

func setColorGray(imgsrc image.Image, channel string) image.Image {
	maxX, maxY := imgsrc.Bounds().Max.X, imgsrc.Bounds().Max.Y
	imgdst := image.NewGray(imgsrc.Bounds())
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			r, g, b := uint32(0), uint32(0), uint32(0)
			switch channel {
			case "red":
				r, _, _, _ = imgsrc.At(x, y).RGBA()
			case "green":
				_, g, _, _ = imgsrc.At(x, y).RGBA()
			case "blue":
				_, _, b, _ = imgsrc.At(x, y).RGBA()
			}
			imgdst.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return imgdst
}
