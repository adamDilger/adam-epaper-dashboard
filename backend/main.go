package main

import (
	"epaper-dashboard/bom"
	"epaper-dashboard/images/bomsummary"
	eastercountdown "epaper-dashboard/images/easter"
	"epaper-dashboard/processing"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	WIDTH  = 800
	HEIGHT = 480
)

func main() {
	imageIndex := 0

	// make a new http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving image")

		var image Image

		if imageIndex == 0 {
			image = weatherSummary()
			imageIndex = 1
		} else {
			image = easterCountdown()
			imageIndex = 0
		}

		headers := []uint8{
			1,                     // format version
			image.durationMinutes, // duration in minutes to display
		}

		response := append(headers, image.data...)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response)))
		w.Write(response)
	})

	port := "8000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Image struct {
	durationMinutes uint8
	data            []byte
}

func weatherSummary() Image {
	a, err := bom.GetBomSummary()
	if err != nil {
		panic(err)
	}

	image := bomsummary.BomSummaryImage(WIDTH, HEIGHT, a)
	data := processing.ConvertContextToBoolArray(image)
	bytesRLE := processing.ConvertBoolArrayToBytesRLE(data)

	return Image{
		durationMinutes: 5,
		data:            bytesRLE,
	}
}

func easterCountdown() Image {
	image := eastercountdown.EasterCountdownImage(WIDTH, HEIGHT, time.Now())
	data := processing.ConvertContextToBoolArray(image)
	bytesRLE := processing.ConvertBoolArrayToBytesRLE(data)

	return Image{
		durationMinutes: 2,
		data:            bytesRLE,
	}
}
