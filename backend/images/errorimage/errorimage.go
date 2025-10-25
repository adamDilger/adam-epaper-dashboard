package errorimage

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed helvetica.ttf
var helvetica embed.FS

func ErrorImage(WIDTH, HEIGHT int, now time.Time) image.Image {
	fonts := sync.OnceValue(loadFonts)()

	im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	dc := gg.NewContextForImage(im)

	dc.SetColor(color.Black)
	dc.SetFontFace(fonts.helvetica.medium)

	dc.DrawStringAnchored("Error :'|", float64(WIDTH)/2, float64(HEIGHT)/2, 0.5, 0.5)

	dc.SetFontFace(fonts.helvetica.extrasmall)
	dc.DrawStringAnchored(
		fmt.Sprintf("Time: %s", now.Format("2006-01-02 15:04")),
		float64(WIDTH)/2,
		float64(HEIGHT)/2+80,
		0.5,
		0.5,
	)

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
