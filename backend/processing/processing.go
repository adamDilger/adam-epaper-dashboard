package processing

import (
	"fmt"

	"github.com/fogleman/gg"
)

func ConvertContextToBoolArray(dc *gg.Context) [][]bool {
	img := dc.Image()

	height := img.Bounds().Max.Y
	width := img.Bounds().Max.X

	// create a 2D array of bools of false values
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width, width)
	}

	for y := range height {
		for x := range width {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 0x7777 {
			} else {
				data[y][x] = true
			}
		}
	}

	return data
}

func ConvertBoolArrayToBytes(data [][]bool) []byte {
	HEIGHT := len(data)
	WIDTH := len(data[0])

	bytes := make([]byte, WIDTH*HEIGHT/8)
	for y := range HEIGHT {
		for x := range WIDTH / 8 {
			for i := range 8 {
				if data[y][x*8+i] {
					bytes[y*(WIDTH/8)+x] |= 1 << (7 - i)
				}
			}
		}
	}

	return bytes
}

func ConvertBoolArrayToBytesRLE(data [][]bool) []uint8 {
	HEIGHT := len(data)
	WIDTH := len(data[0])

	// 8  -- colour of pixel
	// 7  -- count
	// 6  -- count
	// 5  -- count
	// 4  -- count
	// 3  -- count
	// 2  -- count
	// 1  -- count

	fmt.Println(HEIGHT)
	fmt.Println(WIDTH)

	bytes := []uint8{}

	for y := range HEIGHT {
		isBlack := false
		var count uint8 = 1
		x := 0

		for {
			if x >= WIDTH {
				break
			}

			isBlack = data[y][x]

			peek := x + 1
			for peek < WIDTH && count < 127 {
				if data[y][peek] != isBlack {
					break
				}
				peek++
				count++
			}

			x = peek

			var value uint8 = 0
			if isBlack {
				value = 128
			}

			value |= count
			bytes = append(bytes, value)

			count = 1
		}
	}

	return bytes
}
