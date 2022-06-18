package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/db"
	"github.com/margostino/job-pulse/scrapper"
	"github.com/margostino/job-pulse/utils"
	"strings"
)

const (
	SearchPosition = "software engineer"
	SearchLocation = "stockholm"
)

func main() {
	config := configuration.GetConfig()

	var index = 0
	var isEnd = false
	var factor = config.App.ScanFactor
	var baseUrl = config.JobSite.BaseUrl
	var documents = make([]interface{}, 0)
	var fullCardSelector = config.JobSite.FullCardSelector
	var cardInfoSelector = config.JobSite.CardInfoSelector

	var collection = db.ConnectCollection(config.Mongo)

	partialUrl := getPartialUrlFromJobSite(baseUrl)
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
					jobPost := db.BuildJobPost(jobTextParts, sanitizedUrl, rawPostDate, postDate)
					result := db.ConditionalInsert(collection, jobPost)
					documents = append(documents, result)
				}
			}
		}
	}

	db.InsertBatch(collection, documents)
}

func getPartialUrlFromJobSite(baseUrl string) string {
	params := fmt.Sprintf("?keywords=%s&location=%s", SearchPosition, SearchLocation)
	return fmt.Sprintf("%s%s", baseUrl, params)
}
