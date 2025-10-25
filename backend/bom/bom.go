package bom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type RainData struct {
	HourStart        int
	HourEnd          int
	RainfallMills    int
	ChancePercentage int
}

type BomSummary struct {
	LocationName string
	CurrentTemp  string
	TodaysMax    string
	Humidity     string
	Summary      string
	IconName     string
	Rain         []RainData
}

func getBomSummaryJson() (io.ReadCloser, io.ReadCloser, error) {
	weatherUrl := "https://api.bom.gov.au/apikey/v1/observations/latest/40913/atm/surf_air?include_qc_results=false"
	forecastUrl := "https://api.bom.gov.au/apikey/v1/forecasts/texts?aac=QLD_PW015&aac=QLD_FW015&aac=QLD_MW013&aac=QLD_FA001&aac=QLD_ME001&aac=QLD_PT001&timezone=Australia%2FBrisbane"

	resp, err := http.NewRequest(http.MethodGet, weatherUrl, nil)
	resp.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")

	if err != nil {
		return nil, nil, err
	}

	weatherResponse, err := http.DefaultClient.Do(resp)

	if err != nil {
		return nil, nil, err
	}

	forecastResp, err := http.NewRequest(http.MethodGet, forecastUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	forecastResp.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")

	forecastResponse, err := http.DefaultClient.Do(forecastResp)

	if err != nil {
		return nil, nil, err
	}

	println("Fetched BOM data")

	return weatherResponse.Body, forecastResponse.Body, nil
}

func toSafeTempFloat(temp float64) string {
	if temp == 0 {
		return "n/a"
	}

	return fmt.Sprintf("%.1f", temp)
}

func GetBomSummaryTest(path string) (BomSummary, error) {
	f, err := os.Open(path + "test_weather.json")
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to parse file: %v", err)
	}

	ff, err := os.Open(path + "test_forecast.json")
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to parse forecast file: %v", err)
	}

	bytesWeather, err := io.ReadAll(f)
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to read file: %v", err)
	}

	bytesForecast, err := io.ReadAll(ff)
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to read forecast file: %v", err)
	}

	return parseJson(strings.NewReader(string(bytesWeather)), strings.NewReader(string(bytesForecast)))
}

var lock = sync.Mutex{}
var lastFetchTime time.Time
var lastBomSummary BomSummary

const fetchInterval = 10 * time.Minute

func GetBomSummary() (BomSummary, error) {
	if time.Since(lastFetchTime) < fetchInterval {
		log.Printf("Using cached BOM data, since: %v \n", time.Since(lastFetchTime))
		return lastBomSummary, nil
	}

	lock.Lock()
	defer lock.Unlock()

	lastFetchTime = time.Now()

	weatherBody, forecastBody, err := getBomSummaryJson()
	if err != nil {
		lastBomSummary = BomSummary{}
		err = fmt.Errorf("failed to request BOM data: %v", err)
		log.Println(err)
		return lastBomSummary, err
	}

	defer weatherBody.Close()
	defer forecastBody.Close()

	lastBomSummary, err = parseJson(weatherBody, forecastBody)
	if err != nil {
		lastBomSummary = BomSummary{}
		err = fmt.Errorf("failed to parse BOM data: %v", err)
		log.Println(err)
		return lastBomSummary, err
	}

	return lastBomSummary, nil
}

func parseJson(weatherReader io.Reader, forecastReader io.Reader) (BomSummary, error) {
	var weather WeatherResponse
	var forecast ForecastResponse

	err := json.NewDecoder(weatherReader).Decode(&weather)
	if err != nil {
		return BomSummary{}, errors.New("failed to parse weather json: %v" + err.Error())
	}

	err = json.NewDecoder(forecastReader).Decode(&forecast)
	if err != nil {
		return BomSummary{}, errors.New("failed to parse forecast json: %v" + err.Error())
	}

	locationName := "Greenslopes" // TODO: get from response
	// locationName = locationName[0:strings.Index(locationName, "Weather")]

	currentTemp := toSafeTempFloat(weather.Observation.Temperature.DryBulb1MinCel)
	todaysMax := toSafeTempFloat(weather.Observation.Temperature.DryBulbMaxCel)

	summary := forecast.Forecast.Daily[0].Atmospheric.SurfaceAir.Weather.PrecisText

	// iconHref := "" // TODO: doc.Find(".forecasts .forecast-summary dd.image img").First().AttrOr("src", "")
	iconName := "" // TODO: iconHref[strings.LastIndex(iconHref, "/")+1:]

	humidity := fmt.Sprintf("%1.f%%", weather.Observation.Temperature.RelativeHumidityPercent)

	result := BomSummary{
		LocationName: locationName,
		CurrentTemp:  currentTemp,
		TodaysMax:    todaysMax,
		Humidity:     humidity,
		Summary:      summary,
		IconName:     iconName,
		Rain:         []RainData{},
	}

	return result, nil
}
