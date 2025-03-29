package main

import (
	"epaper-dashboard/bom"
	"epaper-dashboard/images/bomsummary"
	"epaper-dashboard/processing"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
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

	image := bomsummary.BomSummaryImage(WIDTH, HEIGHT, a)
	data := processing.ConvertContextToBoolArray(image)
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
