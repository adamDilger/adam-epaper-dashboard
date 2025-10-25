package errorimage_test

import (
	"epaper-dashboard/images/errorimage"
	"epaper-dashboard/processing"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"
)

func TestErrorImage(t *testing.T) {
	// not a real test, just writes the image to a file
	// so we can see what it looks like

	image := errorimage.ErrorImage(800, 480, time.Now())
	data := processing.ConvertContextToBoolArray(image)
	writeImage(data)
}

func writeImage(data [][]bool) {
	HEIGHT := len(data)
	WIDTH := len(data[0])

	i := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	for y := range HEIGHT {
		for x := range WIDTH {
			if data[y][x] {
				i.Set(x, y, color.Black)
			} else {
				i.Set(x, y, color.White)
			}
		}
	}

	o, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}

	err = png.Encode(o, i)
	if err != nil {
		panic(err)
	}
}
