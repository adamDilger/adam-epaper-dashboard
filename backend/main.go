package main

import (
	"epaper-dashboard/bom"
	"epaper-dashboard/processing"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

const (
	WIDTH  = 800
	HEIGHT = 480
)

func main() {
	a, err := bom.GetBomSummaryTest("./test.html")
	if err != nil {
		panic(err)
	}

	fmt.Printf("LocationName : %s\n", a.LocationName)
	fmt.Printf("Current temp: %s\n", a.CurrentTemp)
	fmt.Printf("Todays max: %s\n", a.TodaysMax)
	fmt.Printf("Summary: %s\n", a.Summary)
	fmt.Printf("IconName: %s\n", a.IconName)

	fonts := loadFonts()

	// create an image
	im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	dc := gg.NewContextForImage(im)

	// blank background, not needed, just use alpha value instead
	// dc.SetColor(color.White)
	// dc.DrawRectangle(0, 0, WIDTH, HEIGHT)
	// dc.Fill()

	dc.SetColor(color.Black)
	dc.DrawCircle(80+100, 90+100, 100)
	dc.SetLineWidth(4)
	dc.Stroke()

	dc.SetColor(color.Black)

	dc.SetFontFace(fonts.helvetica.medium)
	dc.DrawStringAnchored(a.LocationName, 30, 400, 0, 0.5)

	dc.SetFontFace(fonts.helvetica.extralarge)
	dc.DrawStringAnchored(a.CurrentTemp, 750, HEIGHT/2, 1, 0.5)

	dc.SetFontFace(fonts.helvetica.small)
	dc.DrawStringAnchored(a.Summary, 750, 400, 1, 0.5)

	dc.SetFontFace(fonts.helvetica.extrasmall)
	w, _ := dc.MeasureString("Max")
	dc.DrawStringAnchored("Max", 750, 50, 1, 0)

	dc.SetFontFace(fonts.helvetica.small)
	dc.DrawStringAnchored(a.TodaysMax, 745-w, 50, 1, 0)

	// load image from file
	iconImage, err := gg.LoadImage(fmt.Sprintf("./images/%s", a.IconName))
	if err != nil {
		panic(err)
	}
	dc.DrawImage(iconImage, 30, 30)

	data := processing.ConvertContextToBoolArray(dc)
	// MakeTestImage(data)

	// bytes := processing.ConvertBoolArrayToBytes(data)
	// fmt.Printf("size of bytes: %d\n", len(bytes))
	// MakeTestImageBytes(bytes)

	bytesRLE := processing.ConvertBoolArrayToBytesRLE(data)
	fmt.Printf("size of bytes RLE: %d\n", len(bytesRLE))
	MakeTestImageBytesRLE(bytesRLE)

	// make a new http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving image")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytesRLE)))
		w.Write(bytesRLE)
	})
	fmt.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type fonts struct {
	helvetica struct {
		extrasmall, small, medium, large, extralarge font.Face
	}
}

func loadFonts() fonts {
	f := fonts{}

	run := func(size float64, dest *font.Face) {
		font, e := gg.LoadFontFace("./helvetica.ttf", size)
		if e != nil {
			panic(e)
		}
		*dest = font
	}

	run(18, &f.helvetica.extrasmall)
	run(30, &f.helvetica.small)
	run(40, &f.helvetica.medium)
	run(60, &f.helvetica.large)
	run(148, &f.helvetica.extralarge)

	return f
}

func MakeTestImage(data [HEIGHT][WIDTH]bool) {
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

func MakeTestImageBytes(data []byte) {
	i := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	for y := range HEIGHT {
		for x := range WIDTH {
			pos := y*(WIDTH/8) + x/8
			bit := x % 8

			if data[pos]&(1<<uint(7-bit)) != 0 {
				i.Set(x, y, color.Black)
			} else {
				i.Set(x, y, color.White)
			}
		}
	}

	o, err := os.Create("out_bytes.png")
	if err != nil {
		panic(err)
	}

	err = png.Encode(o, i)
	if err != nil {
		panic(err)
	}
}

func MakeTestImageBytesRLE(data []uint8) {
	i := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	y := 0
	x := 0
	for _, b := range data {
		isBlack := b&0b10000000 != 0
		count := b & 0b01111111

		for range count {
			if isBlack {
				i.Set(x, y, color.Black)
			} else {
				i.Set(x, y, color.White)
			}

			x++
		}

		if x >= WIDTH {
			x = 0
			y++
		}
	}

	o, err := os.Create("out_bytes_rle.png")
	if err != nil {
		panic(err)
	}

	err = png.Encode(o, i)
	if err != nil {
		panic(err)
	}
}
