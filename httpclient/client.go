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

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	for i := 1; i <= c.maxRetry; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && (resp.StatusCode >= 200 && resp.StatusCode < 400) {
			return
		}
	}

	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	log.Printf(
		"failed to perform [%s] on [%s] error [%v] status code [%d]\n",
		req.Method,
		req.URL,
		err,
		statusCode,
	)

	return
}

func New(timeout, maxRetry int) *Client {
	timeoutDuration := time.Duration(time.Duration(timeout) * time.Second)
	return &Client{
		httpClient: &http.Client{Timeout: timeoutDuration},
		maxRetry:   maxRetry,
	}
}
