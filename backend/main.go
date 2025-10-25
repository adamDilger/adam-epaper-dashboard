package main

import (
	"epaper-dashboard/bom"
	"epaper-dashboard/images/bomsummary"
	"epaper-dashboard/images/errorimage"
	"epaper-dashboard/processing"
	"fmt"
	"image"
	"image/png"
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
	// make a new http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/octet-stream" {
			// coming from the browser
			// send back the image

			var m image.Image
			bomData, err := bom.GetBomSummary()
			if err != nil {
				log.Println("Error getting BOM summary:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error fetching data: " + err.Error()))
				return
			}

			m = bomsummary.BomSummaryImage(WIDTH, HEIGHT, bomData)

			w.Header().Set("Content-Type", "image/png")
			err = png.Encode(w, m)
			if err != nil {
				log.Println("Error encoding image:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error encoding image: " + err.Error()))
			}

			return
		}

		log.Println("Serving image")

		image := weatherSummary()
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
		image := errorimage.ErrorImage(WIDTH, HEIGHT, time.Now())
		data := processing.ConvertContextToBoolArray(image)

		return Image{
			durationMinutes: 5,
			data:            processing.ConvertBoolArrayToBytesRLE(data),
		}
	}

	image := bomsummary.BomSummaryImage(WIDTH, HEIGHT, a)
	data := processing.ConvertContextToBoolArray(image)
	bytesRLE := processing.ConvertBoolArrayToBytesRLE(data)

	// simple hax to reduce the amount of night refreshes
	// TODO: calculate a better duration to ensure the screen always "wakes up" at say 5am
	var duration uint8 = 5
	if time.Now().Hour() >= 0 && time.Now().Hour() < 4 {
		duration = 60
	}

	return Image{
		durationMinutes: duration,
		data:            bytesRLE,
	}
}

/*
func easterCountdown() Image {
	image := eastercountdown.EasterCountdownImage(WIDTH, HEIGHT, time.Now())
	data := processing.ConvertContextToBoolArray(image)
	bytesRLE := processing.ConvertBoolArrayToBytesRLE(data)

	return Image{
		durationMinutes: 2,
		data:            bytesRLE,
	}
}

func isBeforeEaster() bool {
	loc, err := time.LoadLocation("Australia/Brisbane")
	if err != nil {
		panic(err)
	}
	easterDate := time.Date(2025, 4, 21, 0, 0, 0, 0, loc)

	return time.Now().Before(easterDate)
}
*/
