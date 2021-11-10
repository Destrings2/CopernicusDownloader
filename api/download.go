package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const ODATA_ENDPOINT = "/odata/v1/Products('%s')/"

type QuotaExceededError struct {
	Err string
}

func (e QuotaExceededError) Error() string {
	return e.Err
}

type OtherError struct {
	Err string
}

func (e OtherError) Error() string {
	return e.Err
}

type CantRequestOnlineError struct {
	Err string
}

func (e CantRequestOnlineError) Error() string {
	return e.Err
}

type TooManyRequestsError struct {
	Err string
}

func (e TooManyRequestsError) Error() string {
	return e.Err
}

type DownloadRequest struct {
	Title string
	Path  string
	Uuid  string
}

func (cc *CopernicusClient) OpenDownloadChannel() chan DownloadRequest {
	return make(chan DownloadRequest, 100)
}

func (cc *CopernicusClient) Download(uuid string, path string, filename string) (string, error) {
	url := cc.url + fmt.Sprintf(ODATA_ENDPOINT, uuid) + "$value"
	out, err := os.Create(filepath.Join(path, filename) + ".zip")
	if err != nil {
		return "", OtherError{"Download error " + err.Error()}
	}
	defer out.Close()

	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(cc.Username, cc.password)
	if err != nil {
		panic(err)
	}

	resp, err := cc.client.Do(request)
	if err != nil {
		return "", OtherError{"Download error " + err.Error()}
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", OtherError{"Download error " + err.Error()}
	}

	log.Println(fmt.Sprintf("Downloaded %d bytes at %s in path %s", n, filename, path))

	return out.Name(), nil
}

func (cc *CopernicusClient) Request(uuid string) error {
	url := cc.url + fmt.Sprintf(ODATA_ENDPOINT, uuid) + "$value"

	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(cc.Username, cc.password)

	if err != nil {
		panic(err)
	}

	response, err := cc.client.Do(request)
	if response.StatusCode == 202 {
		return nil
	} else if response.StatusCode == 200 {
		return CantRequestOnlineError{"File is actually online"}
	} else if response.StatusCode == 403 {
		return QuotaExceededError{"Quota exceeded at " + time.Now().Format(time.RFC3339)}
	} else if response.StatusCode == 429 {
		return TooManyRequestsError{"Too many requests " + time.Now().Format(time.RFC3339)}
	} else {
		return OtherError{fmt.Sprintf("Status code %d at  %s", response.StatusCode, time.Now().Format(time.RFC3339))}
	}

	return nil
}
