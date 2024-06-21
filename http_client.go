package feidee

import (
	"net/http"
	"net/url"
	"strings"
)

const httpUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"

func (cli *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", httpUserAgent)

	return cli.httpClient.Do(req)
}

func (cli *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", httpUserAgent)

	return cli.httpClient.Do(req)
}
