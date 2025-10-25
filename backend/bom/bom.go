package bom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	IconCode     int
	Rain         []RainData
}

func makeRequest(url, description string) (io.ReadCloser, error) {
	resp, err := http.NewRequest(http.MethodGet, url, nil)
	resp.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")

	if err != nil {
		return nil, fmt.Errorf("failed to build request for %s: %v", description, err)
	}

	response, err := http.DefaultClient.Do(resp)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response for %s: %d", description, response.StatusCode)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to perform request for %s: %v", description, err)
	}

	return response.Body, nil
}

func getBomSummaryJson() (io.ReadCloser, io.ReadCloser, io.ReadCloser, error) {
	weatherUrl := "https://api.bom.gov.au/apikey/v1/observations/latest/40913/atm/surf_air?include_qc_results=false"
	forecastTextUrl := "https://api.bom.gov.au/apikey/v1/forecasts/texts?aac=QLD_PW015&aac=QLD_FW015&aac=QLD_MW013&aac=QLD_FA001&aac=QLD_ME001&aac=QLD_PT001&timezone=Australia%2FBrisbane"
	forecastDailyUrl := "https://api.bom.gov.au/apikey/v1/forecasts/daily/689/350?timezone=Australia%2FBrisbane"

	var weatherResponse, forecastDailyResponse, forecastTextReponse io.ReadCloser
	var err1, err2, err3 error

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		weatherResponse, err1 = makeRequest(weatherUrl, "Weather")
		wg.Done()
	}()
	go func() {
		forecastTextReponse, err2 = makeRequest(forecastTextUrl, "Forecast Text")
		wg.Done()
	}()
	go func() {
		forecastDailyResponse, err3 = makeRequest(forecastDailyUrl, "Forecast Daily")
		wg.Done()
	}()

	wg.Wait()

	for _, err := range []error{err1, err2, err3} {
		if err != nil {
			return nil, nil, nil, err
		}
	}

	println("Fetched BOM data")

	return weatherResponse, forecastDailyResponse, forecastTextReponse, nil
}

func toSafeTempFloat(temp float64) string {
	if temp == 0 {
		return "n/a"
	}

	return fmt.Sprintf("%.1f°", temp)
}

func GetBomSummaryTest(path string) (BomSummary, error) {
	weatherFile, err := os.Open(path + "test_weather.json")
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to parse file: %v", err)
	}

	forecastTextsFile, err := os.Open(path + "test_forecast_texts.json")
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to parse forecast file: %v", err)
	}

	forecastDailyFile, err := os.Open(path + "test_forecast_daily.json")
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to parse forecast file: %v", err)
	}

	return parseJson(weatherFile, forecastDailyFile, forecastTextsFile)
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

	weatherResponse, forecastDailyResponse, forecastTextResponse, err := getBomSummaryJson()
	if err != nil {
		return lastBomSummary, fmt.Errorf("failed to request BOM data: %v", err)
	}

	defer weatherResponse.Close()
	defer forecastDailyResponse.Close()
	defer forecastTextResponse.Close()

	lastBomSummary, err = parseJson(weatherResponse, forecastDailyResponse, forecastTextResponse)
	if err != nil {
		return lastBomSummary, fmt.Errorf("failed to parse BOM data: %v", err)
	}

	return lastBomSummary, nil
}

func parseJson(weatherReader, forecastDailyReader, forecastTextsReader io.Reader) (BomSummary, error) {
	var weather WeatherResponse
	var forecastTexts ForecastTextsResponse
	var forecastDaily DailyForecastsResponse

	err := json.NewDecoder(weatherReader).Decode(&weather)
	if err != nil {
		return BomSummary{}, errors.New("failed to parse weather json: %v" + err.Error())
	}

	err = json.NewDecoder(forecastDailyReader).Decode(&forecastDaily)
	if err != nil {
		return BomSummary{}, errors.New("failed to parse forecast daily json: %v" + err.Error())
	}

	err = json.NewDecoder(forecastTextsReader).Decode(&forecastTexts)
	if err != nil {
		return BomSummary{}, errors.New("failed to parse forecast texts json: %v" + err.Error())
	}

	locationName := "Greenslopes" // TODO: get from response

	currentTemp := toSafeTempFloat(weather.Observation.Temperature.DryBulb1MinCel)
	todaysMax := toSafeTempFloat(weather.Observation.Temperature.DryBulbMaxCel)

	summary := forecastTexts.Forecast.Daily[0].Atmospheric.SurfaceAir.Weather.PrecisText
	iconCode := forecastDaily.ForecastsDaily.Daily[0].Atmospheric.SurfaceAir.Weather.IconCode
	humidity := fmt.Sprintf("%1.f%%", weather.Observation.Temperature.RelativeHumidityPercent)

	result := BomSummary{
		LocationName: locationName,
		CurrentTemp:  currentTemp,
		TodaysMax:    todaysMax,
		Humidity:     humidity,
		Summary:      summary,
		IconCode:     iconCode,
		Rain:         []RainData{},
	}

	return result, nil
}
