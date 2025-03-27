package processing

import "github.com/fogleman/gg"

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
