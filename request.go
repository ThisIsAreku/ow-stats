package ow_stats

import (
	"net/http"
	"time"
)

type HttpClient struct {
	baseClient *http.Client
	userAgent  string
}

func NewHttpClient() (*HttpClient) {
	return &HttpClient{
		&http.Client{
			Timeout: time.Second * 30,
			Jar:     nil,
		},
		"OW-STATS/1.0",
	}
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	return c.baseClient.Do(req)
}
