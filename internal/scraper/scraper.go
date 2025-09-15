// Package scraper involves all functions related to actually doing http requests and scraping the
// data from the response
package scraper

import (
	"errors"
	"net/http"
)

func Scrape(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", errors.New(err.Error())
	}

	contentType := res.Header.Get("Content-Type")
	return contentType, nil
}
