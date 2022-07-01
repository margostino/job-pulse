package collector

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/job-pulse/domain"
	"github.com/margostino/job-pulse/utils"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

func (a *App) ValidateInput() {
	if len(os.Args) != 3 {
		example := "go run ./collector \"engineer\" \"stockholm\""
		message := fmt.Sprintf("Missing parameters!\nExample: %s", example)
		panic(message)
	}
	a.inputParams.SearchPosition = os.Args[1]
	a.inputParams.SearchLocation = os.Args[2]
}

func (a *App) Start() error {
	stats := &domain.Stats{
		StartTime:     time.Now().UTC(),
		PositionInput: a.inputParams.SearchPosition,
		LocationInput: a.inputParams.SearchLocation,
	}

	var index = 0
	var isEnd = false
	var factor = a.config.App.ScanFactor
	var baseUrl = a.config.JobSite.BaseUrl
	var documents = make([]interface{}, 0)
	var cardInfoSelector = a.config.JobSite.CardInfoSelector

	partialUrl := a.getPartialUrlFromJobSite(baseUrl)

	for ok := true; ok; ok = !isEnd {
		index += factor
		fullUrl := utils.IndexFrom(partialUrl, index)
		webScrapper := a.scrapper.GoPage(fullUrl)
		entries := webScrapper.GetLinkElements()

		if len(entries) == 0 {
			if a.isTimeModeActive(stats) {
				index = 0
				a.db.InsertBatch(documents, stats)
				documents = make([]interface{}, 0)
			} else {
				isEnd = true
			}
		}

		for _, entry := range entries {
			var sanitizedUrl = a.extractUrlFrom(entry)
			card, err := entry.Element(cardInfoSelector)
			if err != nil {
				log.Println(err.Error())
				log.Println(entry.String())
				break
			} else {
				jobText, err := card.Text()
				utils.Check(err)
				jobText = strings.ToLower(jobText)
				jobTextParts := strings.SplitN(jobText, "\n", -1)

				if len(jobTextParts) > 0 {
					rawPostDate := utils.GetRawPostDateOrDefault(4, jobTextParts)
					postDate := calculateJobPostDate(rawPostDate)
					jobPost := buildJobPost(jobTextParts, sanitizedUrl, rawPostDate, postDate)
					jobPost.Latitude, jobPost.Longitude = a.getGeocoding(jobPost.Location)
					result := a.db.GetConditionalDocument(jobPost)
					if result != nil {
						documents = append(documents, result)
					}
				}
			}
		}
	}

	a.db.InsertBatch(documents, stats)

	return nil // TODO tbd
}

func (a *App) Close() {
	a.db.Close()
}

func (a *App) getPartialUrlFromJobSite(baseUrl string) string {
	params := fmt.Sprintf("?keywords=%s&location=%s", a.inputParams.SearchPosition, a.inputParams.SearchLocation)
	return fmt.Sprintf("%s%s", baseUrl, params)
}

func (a *App) isTimeModeActive(stats *domain.Stats) bool {
	return a.config.App.TimeMode && timeDurationSince(stats.StartTime) < a.config.App.TimeModeDuration
}

func (a *App) extractUrlFrom(entry *rod.Element) string {
	var sanitizedUrl string
	var fullCardSelector = a.config.JobSite.FullCardSelector
	fullCard, err := entry.Element(fullCardSelector)
	if err != nil {
		sanitizedUrl = ""
	} else {
		linkSource, err := fullCard.Property("href")
		utils.Check(err)
		sanitizedUrl = utils.SanitizeUrl(linkSource.String())
		utils.Check(err)
	}
	return sanitizedUrl
}

// TODO: separate concerns between Geo and DB modules
func (a *App) getGeocoding(location string) (float64, float64) {
	var latitude float64
	var longitude float64
	geocoding, _ := a.db.FindOneGeoBy(location)
	// TODO: best effort for location similarity string (e.g. Stockholm == Stockholm, Sweden)
	if geocoding == nil {
		log.Printf("New geocoding for location: %s", location)
		newGeocoding := a.geo.Get(location)
		if newGeocoding != nil {
			println(fmt.Sprintf("%v", *newGeocoding))
			validateGeoResponse(*newGeocoding)
			newGeocodingMap := (*newGeocoding).(map[string]interface{})
			latitude = newGeocodingMap["latitude"].(float64)
			longitude = newGeocodingMap["longitude"].(float64)
			a.db.InsertOneGeocoding(location, newGeocoding)
		}
	} else {
		latitude = geocoding["latitude"].(float64)
		longitude = geocoding["longitude"].(float64)
	}
	return latitude, longitude
}

func validateGeoResponse(value interface{}) {
	rt := reflect.TypeOf(value)
	if rt.Kind() == reflect.Slice {
		log.Println("Geocoding response was empty slice")
	}
}
