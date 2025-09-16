// Package scraper involves all functions related to actually doing http requests and scraping the
// data from the response
package scraper

import (
	"errors"
	"net/http"
	"time"

	l "github.com/JesterSe7en/scrapego/internal/logger"
	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

func Scrape(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", errors.New(err.Error())
	}

	contentType := res.Header.Get("Content-Type")
	return contentType, nil
}

func ScrapeWithTimeout(url string, timeout time.Duration) wp.Result {
	res, err := http.Get(url)
	if err != nil {
		return wp.Result{
			Value: nil,
			Err:   err,
		}
	}

	contentType := res.Header.Get("Content-Type")
	l.Debug("contentType = %w", contentType)
	return wp.Result{
		Value: contentType,
		Err:   nil,
	}
}
