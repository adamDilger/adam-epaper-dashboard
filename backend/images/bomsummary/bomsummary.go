package bomsummary

import (
	"embed"
	"epaper-dashboard/bom"
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

//go:embed icons/*.png
var icons embed.FS

func BomSummaryImage(WIDTH, HEIGHT int, a bom.BomSummary) image.Image {
	fonts := sync.OnceValue(loadFonts)()

	im := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	dc := gg.NewContextForImage(im)

	dc.SetColor(color.Black)

	dc.SetFontFace(fonts.helvetica.small)
	dc.DrawStringAnchored(a.LocationName, 30, 420, 0, 0.5)

	dc.SetFontFace(fonts.helvetica.extralarge)
	dc.DrawStringAnchored(a.CurrentTemp, 750, float64(HEIGHT)/2, 1, 0.5)

	if len(a.Summary) > 13 {
		dc.SetFontFace(fonts.helvetica.medium)
		dc.DrawStringAnchored(a.Summary, 750, 430, 1, 0)
	} else {
		dc.SetFontFace(fonts.helvetica.large)
		dc.DrawStringAnchored(a.Summary, 750, 430, 1, 0)
	}

	dc.SetFontFace(fonts.helvetica.medium)
	w, _ := dc.MeasureString(a.TodaysMax)
	dc.DrawStringAnchored(a.TodaysMax, 750, 70, 1, 0)

	dc.SetFontFace(fonts.helvetica.extrasmall)
	dc.DrawStringAnchored("Max", 740-w, 70, 1, 0)

	dc.SetFontFace(fonts.helvetica.medium)
	w, _ = dc.MeasureString(a.Humidity)
	dc.DrawStringAnchored(a.Humidity, 750, 120, 1, 0)

	dc.SetFontFace(fonts.helvetica.extrasmall)
	dc.DrawStringAnchored("Humidity", 740-w, 120, 1, 0)

	iconCode := a.IconCode
	if iconDefinition, ok := IconDefinitionMap[iconCode]; ok {
		iconImageFile, err := icons.Open(fmt.Sprintf("icons/%s.png", iconDefinition.DayIconName))
		if err != nil {
			fmt.Println("Failed to open icon image:", err)
		} else {
			iconImage, _, err := image.Decode(iconImageFile)
			if err != nil {
				fmt.Println("Failed to decode icon image:", err)
			} else {
				dc.DrawImage(iconImage, 10, 10)
			}
		}
	}

	nowString := time.Now().Format("3:04pm 2/1/06")
	dc.SetFontFace(fonts.helvetica.extraextrasmall)
	dc.DrawStringAnchored(nowString, 800, 480, 1, 0)

	// dc.SetRGB(1, 1, 1)
	// dc.SetLineWidth(2)
	// dc.SetFontFace(fonts.helvetica.extrasmall)
	//
	// dc.DrawLine(350, 100, float64(350+(len(a.Rain)*30)), 100)
	// dc.Stroke()
	// dc.DrawLine(350, 50, float64(350+(len(a.Rain)*30)), 50)
	// dc.Stroke()
	//
	// for i, rain := range a.Rain {
	// 	x := float64(i*30) + 350
	// 	dc.DrawLine(x, 100, x, float64(100-(rain.ChancePercentage/2)))
	// 	dc.Stroke()
	//
	// 	dc.DrawStringAnchored(fmt.Sprintf("%d", rain.RainfallMills), x, float64(90-(rain.RainfallMills*10)), 0.5, 0)
	// 	dc.DrawStringAnchored(fmt.Sprintf("%d", rain.HourStart), x, 100, 0.5, 1)
	// }

	return dc.Image()
}

type fonts struct {
	helvetica struct {
		extraextrasmall, extrasmall, small, medium, large, extralarge font.Face
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

	run(14, &f.helvetica.extraextrasmall)
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
