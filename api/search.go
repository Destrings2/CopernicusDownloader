package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	API_ROOT     = "https://scihub.copernicus.eu"
	API_ENDPOINT = "dhus"
)

type CopernicusClient struct {
	Username               string
	password               string
	client                 http.Client
	url                    string
	RequestQuotaExceeded   bool
	RequestQuotaExceededAt time.Time
}

func NewCopernicusClient(username string, password string) *CopernicusClient {
	return &CopernicusClient{
		Username: username,
		password: password,
		client:   http.Client{},
		url:      API_ROOT + "/" + API_ENDPOINT,
	}
}

func (cc *CopernicusClient) Search(query string, start int, rows int) *Feed {
	queryUrl := fmt.Sprintf("/search?start=%d&rows=%d&q=%s", start, rows, url.QueryEscape(query))
	request, err := http.NewRequest("GET", cc.url+queryUrl, nil)
	request.SetBasicAuth(cc.Username, cc.password)

	if err != nil {
		panic(err)
	}

	response, err := cc.client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var feed Feed
	err = xml.Unmarshal(data, &feed)

	if err != nil {
		log.Println(err)
	}

	return &feed
}

func (cc *CopernicusClient) GetEntryByUUID(uuid string) *ODataEntry {
	queryUrl := fmt.Sprintf(ODATA_ENDPOINT, uuid)
	request, err := http.NewRequest("GET", cc.url+queryUrl, nil)
	request.SetBasicAuth(cc.Username, cc.password)

	if err != nil {
		panic(err)
	}

	response, err := cc.client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var entry ODataEntry
	err = xml.Unmarshal(data, &entry)

	if err != nil {
		log.Println(fmt.Sprintf("Error marshaling XML %s, got response code %d", err, response.StatusCode))
		return nil
	}

	return &entry
}
