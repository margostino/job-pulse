package collector

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/db"
	"github.com/margostino/job-pulse/domain"
	"github.com/margostino/job-pulse/geo"
	"github.com/margostino/job-pulse/scrapper"
	"github.com/margostino/job-pulse/utils"
	"log"
	"os"
	"strings"
	"time"
)

type InputParams struct {
	searchPosition string
	searchLocation string
}

type Collector struct {
	db          *db.Connection
	geo         *geo.Connection
	config      *configuration.Configuration
	inputParams *InputParams
}

func (c *Collector) Start() error {
	if len(os.Args) != 3 {
		example := "go run ./collector \"engineer\" \"stockholm\""
		message := fmt.Sprintf("Missing parameters!\nExample: %s", example)
		panic(message)
	}
	c.inputParams.searchPosition = os.Args[1]
	c.inputParams.searchLocation = os.Args[2]

	var latitude float64
	var longitude float64
	var index = 0
	var isEnd = false
	var factor = c.config.App.ScanFactor
	var baseUrl = c.config.JobSite.BaseUrl
	var documents = make([]interface{}, 0)
	var fullCardSelector = c.config.JobSite.FullCardSelector
	var cardInfoSelector = c.config.JobSite.CardInfoSelector

	partialUrl := c.getPartialUrlFromJobSite(baseUrl)
	browser := rod.New().MustConnect()
	defer browser.Close()

	for ok := true; ok; ok = !isEnd {
		index += factor
		fullUrl := utils.IndexFrom(partialUrl, index)
		page := browser.MustPage(fullUrl).MustWaitLoad()
		entries, err := page.Elements("li")
		utils.Check(err)

		if len(entries) == 0 {
			isEnd = true
			//index = 0
		}

		for _, entry := range entries {
			var sanitizedUrl string
			fullCard, err := entry.Element(fullCardSelector)
			if err != nil {
				sanitizedUrl = ""
			} else {
				linkSource, err := fullCard.Property("href")
				utils.Check(err)
				sanitizedUrl = utils.SanitizeUrl(linkSource.String())
				utils.Check(err)
			}

			card, err := entry.Element(cardInfoSelector)
			if err != nil {
				println(err.Error())
				break
			} else {
				jobText, err := card.Text()
				utils.Check(err)
				jobText = strings.ToLower(jobText)
				jobTextParts := strings.SplitN(jobText, "\n", -1)

				if len(jobTextParts) > 0 {
					rawPostDate := utils.GetRawPostDateOrDefault(4, jobTextParts)
					postDate := scrapper.CalculateJobPostDate(rawPostDate)
					jobPost := buildJobPost(jobTextParts, sanitizedUrl, rawPostDate, postDate)
					geocoding, _ := c.db.FindOneGeoBy(jobPost.Location)
					// TODO: best effort for location similarity string (e.g. Stockholm == Stockholm, Sweden)
					if geocoding == nil {
						log.Printf("New geocoding for %s", jobPost.Location)
						newGeocoding := c.geo.Get(jobPost.Location)
						if newGeocoding != nil {
							newGeocodingMap := (*newGeocoding).(map[string]interface{})
							latitude = newGeocodingMap["latitude"].(float64)
							longitude = newGeocodingMap["longitude"].(float64)
							c.db.InsertOneGeocoding(jobPost.Location, newGeocoding)
						}
					} else {
						latitude = geocoding["latitude"].(float64)
						longitude = geocoding["longitude"].(float64)
					}
					jobPost.Latitude = latitude
					jobPost.Longitude = longitude
					result := c.db.GetConditionalDocument(jobPost)
					if result != nil {
						documents = append(documents, result)
					}
				}
			}
		}
	}

	c.db.InsertBatch(documents)

	return nil // TODO tbd
}

func (c *Collector) Close() {
	c.db.Close()
}

func (c *Collector) getPartialUrlFromJobSite(baseUrl string) string {
	params := fmt.Sprintf("?keywords=%s&location=%s", c.inputParams.searchPosition, c.inputParams.searchLocation)
	return fmt.Sprintf("%s%s", baseUrl, params)
}

func buildJobPost(jobTextParts []string, link string, rawPostDate string, postDate time.Time) *domain.JobPost {
	return &domain.JobPost{
		Position:    utils.GetOrDefault(0, jobTextParts),
		Company:     utils.GetOrDefault(1, jobTextParts),
		Location:    utils.GetOrDefault(2, jobTextParts),
		Benefit:     utils.GetOrDefault(3, jobTextParts),
		Link:        link,
		RawPostDate: rawPostDate,
		PostDate:    postDate,
	}
}

func initInputParams() *InputParams {
	return &InputParams{
		searchPosition: "",
		searchLocation: "",
	}
}
