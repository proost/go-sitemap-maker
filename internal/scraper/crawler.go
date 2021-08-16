package scraper

import (
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

var client = http.Client{Timeout: time.Second * 1}

func GetUrlContent(url string) (string, error) {
	var (
		err error
		content []byte
		resp *http.Response
	)

	// get content of url
	if resp, err = client.Get(url); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		return "", err
	}

	// Read the body of the HTTP response
	if content, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	return html.UnescapeString(string(content)), nil
}
