package bom

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	html, err := getBomSummaryHtml()
	if err != nil {
		lastBomSummary = BomSummary{}
		err = fmt.Errorf("failed to request BOM data: %v", err)
		log.Println(err)
		return lastBomSummary, err
	}

	lastBomSummary, err = parseHtml(html)
	if err != nil {
		lastBomSummary = BomSummary{}
		err = fmt.Errorf("failed to parse BOM data: %v", err)
		log.Println(err)
		return lastBomSummary, err
	}

	return lastBomSummary, nil
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

	/*
		doc.Find(".pme table tbody tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if i == 0 {
				return true // Skip the first row
			}

			time := s.Find("td.time").Text()
			// regex for this string: 7:00 am - 10:00 am
			rg := regexp.MustCompile(`(\d+):\d+\s([ap]m)\s-\s(\d+):\d+\s([ap]m)`)
			matches := rg.FindStringSubmatch(time)

			rd := RainData{}

			fmt.Printf("matches: %v\n", matches)
			if v, err := strconv.Atoi(matches[1]); err != nil {
				log.Printf("failed to parse hour start: %v", err)
			} else {
				rd.HourStart = v
			}

			if matches[2] == "pm" {
				rd.HourStart += 12
			}

			if v, err := strconv.Atoi(matches[3]); err != nil {
				log.Printf("failed to parse hour end: %v", err)
			} else {
				rd.HourEnd = v
			}

			if matches[4] == "pm" {
				rd.HourEnd += 12
			}

			rainfallString := s.Find(".amt").Text()
			if strings.Contains(rainfallString, " - ") {
				rainfallString = strings.TrimSpace(strings.Split(rainfallString, " - ")[0])
				rainfallString = strings.ReplaceAll(rainfallString, " mm", "")
				if v, err := strconv.Atoi(rainfallString); err != nil {
					log.Printf("failed to parse rainfall: %v", err)
				} else {
					rd.RainfallMills = v
				}
			} else {
				rainfallString = strings.ReplaceAll(rainfallString, " mm", "")
				if v, err := strconv.Atoi(rainfallString); err != nil {
					log.Printf("failed to parse rainfall: %v", err)
				} else {
					rd.RainfallMills = v
				}
			}

			chance := s.Find(".coaf").Text()
			chance = strings.TrimSpace(strings.ReplaceAll(chance, "%", ""))
			if v, err := strconv.Atoi(chance); err != nil {
				log.Printf("failed to parse chance: %v", err)
				return true
			} else {
				rd.ChancePercentage = v
			}

			result.Rain = append(result.Rain, rd)
			return true
		})
	*/

	return result, nil
}
