package httpclient

import (
	"log"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	maxRetry   int
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	for i := 1; i <= c.maxRetry; i++ {
		resp, err = c.httpClient.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
	}

	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}
	log.Printf("failed to get [%s] error [%v] status code [%d]\n", url, err, statusCode)

	return
}

func New(timeout, maxRetry int) *Client {
	timeoutDuration := time.Duration(time.Duration(timeout) * time.Second)
	return &Client{
		httpClient: &http.Client{Timeout: timeoutDuration},
		maxRetry:   maxRetry,
	}
}
