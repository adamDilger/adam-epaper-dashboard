package main

import (
	"epaper-dashboard/bom"
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
	dc.DrawStringAnchored("Stones Corner", 30, 400, 0, 0.5)

	dc.SetFontFace(fonts.helvetica.extralarge)
	dc.DrawStringAnchored(a.CurrentTemp, 750, HEIGHT/2, 1, 0.5)

	dc.SetFontFace(fonts.helvetica.small)
	dc.DrawStringAnchored(a.Summary, 750, 400, 1, 0.5)

	// load image from file
	iconImage, err := gg.LoadImage(fmt.Sprintf("./images/%s", a.IconName))
	if err != nil {
		panic(err)
	}
	dc.DrawImage(iconImage, 30, 30)

	data := BuildDataArray(dc)
	bytes := BuildByteArray(data)

	// MakeTestImage(data)
	MakeTestImageBytes(bytes[:])

	// make a new http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving image")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
		w.Write(bytes[:])
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

func PrintTestImage(data [HEIGHT][WIDTH]bool) {
	for y := 0; y < HEIGHT; y++ {
		fmt.Print("|")

		for x := 0; x < WIDTH; x++ {
			if data[y][x] {
				fmt.Print("X")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
}

func BuildDataArray(dc *gg.Context) [HEIGHT][WIDTH]bool {
	data := [HEIGHT][WIDTH]bool{}
	img := dc.Image()

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 0x7777 {
			} else {
				data[y][x] = true
			}
		}
	}

	return data
}

func BuildByteArray(data [HEIGHT][WIDTH]bool) [WIDTH * HEIGHT / 8]byte {
	bytes := [WIDTH * HEIGHT / 8]byte{}
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH/8; x++ {
			for i := 0; i < 8; i++ {
				if data[y][x*8+i] {
					bytes[y*(WIDTH/8)+x] |= 1 << (7 - i)
				}
			}
		}
	}

	return bytes
}

func MakeTestImage(data [HEIGHT][WIDTH]bool) {
	i := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
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

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
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
