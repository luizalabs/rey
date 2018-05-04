package checker

import (
	"github.com/luizalabs/rey/component"
	"github.com/luizalabs/rey/httpclient"
	"github.com/luizalabs/rey/status"
)

const statusDisruption = 500

type Checker struct {
	maxRetry int
	timeout  int
}

func (c *Checker) Check(comp *component.Component) (*status.Status, error) {
	var currentStatus int
	var currentDetail string

	cli := httpclient.New(c.timeout, c.maxRetry)
	resp, err := cli.Get(comp.HCEndpoint)
	if err != nil {
		currentStatus = statusDisruption
		currentDetail = err.Error()
	} else {
		currentStatus = resp.StatusCode
		currentDetail = resp.Status
	}

	st := &status.Status{
		Component: comp.Name,
		Details:   currentDetail,
		Status:    currentStatus,
	}

	return st, nil
}

func New(timeout, maxRetry int) *Checker {
	return &Checker{maxRetry: maxRetry, timeout: timeout}
}
