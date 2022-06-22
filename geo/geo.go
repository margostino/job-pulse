package geo

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Connection struct {
	Client    *http.Client
	AccessKey string
	BaseUrl   string
}

func Connect(config *configuration.Configuration) *Connection {
	client := &http.Client{
		Timeout: time.Second * time.Duration(config.Geo.Timeout),
	}
	return &Connection{
		Client:    client,
		AccessKey: config.Geo.AccessKey,
		BaseUrl:   config.Geo.BaseUrl,
	}
}

func (c *Connection) Get(query string) *interface{} {
	var geocoding interface{}
	request := c.getRequest(query)
	response, err := c.Client.Do(request)
	utils.Check(err)

	body, err := ioutil.ReadAll(response.Body)
	utils.Check(err)
	response.Body.Close()

	if err := json.Unmarshal(body, &geocoding); err != nil {
		log.Fatal(err)
	}
	// TODO: validate and check the following before returning
	rawData := geocoding.(map[string]interface{})["data"]

	if rawData != nil {
		data := rawData.([]interface{})
		if len(data) > 0 {
			geocodingData := geocoding.(map[string]interface{})["data"].([]interface{})[0]
			return &geocodingData
		}
	}
	log.Printf("No geocoding results for query %s", query)
	return nil
}

func (c *Connection) getRequest(query string) *http.Request {
	request, err := http.NewRequest("GET", c.BaseUrl, nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "application/json")
	q := request.URL.Query()
	q.Add("access_key", c.AccessKey)
	q.Add("query", query)
	q.Add("output", "json")
	q.Add("limit", "1")
	request.URL.RawQuery = q.Encode()
	return request
}
