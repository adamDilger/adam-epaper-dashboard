package bom

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RainData struct {
	Time     string
	Rainfall string
	Chance   string
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

func getBomSummaryHtml() (string, error) {
	url := "http://www.bom.gov.au/places/qld/greenslopes"

	resp, err := http.NewRequest(http.MethodGet, url, nil)
	resp.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", err
	}

	return html, nil
}

func toSafeTemp(temp string) string {
	if temp == "" {
		return "n/a"
	}
	return strings.ReplaceAll(temp, " °C", "°")
}

func GetBomSummaryTest(path string) (BomSummary, error) {
	f, err := os.Open(path)
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to parse file: %v", err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to read file: %v", err)
	}

	return parseHtml(string(bytes))
}

func GetBomSummary() (BomSummary, error) {
	html, err := getBomSummaryHtml()
	if err != nil {
		return BomSummary{}, fmt.Errorf("failed to request BOM data: %v", err)
	}

	return parseHtml(html)
}

func parseHtml(html string) (BomSummary, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return BomSummary{}, errors.New("failed to parse HTML")
	}

	locationName := doc.Find("h1").First().Text()
	locationName = locationName[0:strings.Index(locationName, "Weather")]

	currentTemp := toSafeTemp(doc.Find("li.airT").First().Text())
	todaysMax := toSafeTemp(doc.Find("dd.max").First().Text())

	summary := doc.Find(".forecasts .forecast-summary dd.summary").First().Text()

	iconHref := doc.Find(".forecasts .forecast-summary dd.image img").First().AttrOr("src", "")
	iconName := iconHref[strings.LastIndex(iconHref, "/")+1:]

	humidity := doc.Find("#summary-1 td").First().Text()

	result := BomSummary{
		LocationName: locationName,
		CurrentTemp:  currentTemp,
		TodaysMax:    todaysMax,
		Humidity:     humidity,
		Summary:      summary,
		IconName:     iconName,
		Rain:         []RainData{},
	}

	doc.Find(".pme table tbody tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			return true // Skip the first row
		}

		time := strings.TrimSpace(s.Find(".time").Text())
		time = strings.ReplaceAll(strings.ReplaceAll(time, " :00", ""), " (.m)", "$1")
		rainfall := strings.TrimSpace(strings.ReplaceAll(s.Find(".amt").Text(), " ", ""))
		chance := strings.TrimSpace(s.Find(".coaf").Text())

		if time == "" || rainfall == "" || chance == "" {
			return true
		}

		result.Rain = append(result.Rain, RainData{
			Time:     time,
			Rainfall: rainfall,
			Chance:   chance,
		})

		return true
	})

	return result, nil
}
