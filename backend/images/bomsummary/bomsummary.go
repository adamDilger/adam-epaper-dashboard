package bomsummary

import (
	"embed"
	"epaper-dashboard/bom"
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed helvetica.ttf
var helvetica embed.FS

//go:embed icons/*.png
var icons embed.FS

func BomSummaryImage(WIDTH, HEIGHT int, a bom.BomSummary) image.Image {
	fonts := loadFonts()

	im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	dc := gg.NewContextForImage(im)

	dc.SetColor(color.Black)

	dc.SetFontFace(fonts.helvetica.medium)
	dc.DrawStringAnchored(a.LocationName, 30, 400, 0, 0.5)

	dc.SetFontFace(fonts.helvetica.extralarge)
	dc.DrawStringAnchored(a.CurrentTemp, 750, float64(HEIGHT)/2, 1, 0.5)

	dc.SetFontFace(fonts.helvetica.small)
	dc.DrawStringAnchored(a.Summary, 750, 400, 1, 0.5)

	dc.SetFontFace(fonts.helvetica.small)
	w, _ := dc.MeasureString(a.TodaysMax)
	dc.DrawStringAnchored(a.TodaysMax, 750, 50, 1, 0)

	dc.SetFontFace(fonts.helvetica.extrasmall)
	dc.DrawStringAnchored("Max", 740-w, 50, 1, 0)

	dc.SetFontFace(fonts.helvetica.small)
	w, _ = dc.MeasureString(a.Humidity)
	dc.DrawStringAnchored(a.Humidity, 750, 80, 1, 0)

	dc.SetFontFace(fonts.helvetica.extrasmall)
	dc.DrawStringAnchored("Humidity", 740-w, 80, 1, 0)

	iconImageFile, err := icons.Open(fmt.Sprintf("icons/%s", a.IconName))
	if err != nil {
		panic(err)
	}
	iconImage, _, err := image.Decode(iconImageFile)
	dc.DrawImage(iconImage, 10, 10)

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
	run(148, &f.helvetica.extralarge)

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
