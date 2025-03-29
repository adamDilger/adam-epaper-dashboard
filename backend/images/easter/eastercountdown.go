package eastercountdown

import (
	"embed"
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed bunny.*
var bunnyImage embed.FS

//go:embed helvetica.ttf
var helvetica embed.FS

func EasterCountdownImage(WIDTH, HEIGHT int, now time.Time) image.Image {
	fonts := loadFonts()

	im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	dc := gg.NewContextForImage(im)

	dc.SetColor(color.Black)

	loc, err := time.LoadLocation("Australia/Brisbane")
	if err != nil {
		panic(err)
	}
	easterDate := time.Date(2025, 4, 20, 0, 0, 0, 0, loc)

	// days difference from now
	dayCount := int(easterDate.Sub(now).Hours()/24) + 1
	days := strconv.Itoa(dayCount)

	dc.SetFontFace(fonts.helvetica.large)
	dc.DrawStringAnchored("Easter", 200, 100, 0.5, 1)

	dc.SetFontFace(fonts.helvetica.extralarge)
	_, h := dc.MeasureString(days)
	dc.DrawStringAnchored(days, 200, 240, 0.5, 0.5)

	dc.SetFontFace(fonts.helvetica.small)
	dc.DrawStringAnchored("Days", 200, 240+(h/2)+20, 0.5, 1)

	bunnyImageFile, err := bunnyImage.Open("bunny.png")
	if err != nil {
		panic(err)
	}
	iconImage, _, err := image.Decode(bunnyImageFile)
	dc.DrawImage(iconImage, 800-480, 0)

	return dc.Image()
}

type fonts struct {
	helvetica struct {
		extrasmall, small, medium, large, extralarge font.Face
	}
}

func loadFonts() fonts {
	f := fonts{}

	run := func(size float64, dest *font.Face) {
		font, e := LoadFontFace(size)
		if e != nil {
			panic(e)
		}
		*dest = font
	}

	run(18, &f.helvetica.extrasmall)
	run(30, &f.helvetica.small)
	run(40, &f.helvetica.medium)
	run(60, &f.helvetica.large)
	run(168, &f.helvetica.extralarge)

	return f
}

func LoadFontFace(points float64) (font.Face, error) {
	fontBytes, err := helvetica.ReadFile("helvetica.ttf")
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}
