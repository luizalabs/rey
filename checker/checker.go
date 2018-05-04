package checker

import (
	"net/http"

	"github.com/luizalabs/rey/component"
	"github.com/luizalabs/rey/httpclient"
	"github.com/luizalabs/rey/status"
)

const (
	statusOperational = 100
	statusDisruption  = 500
)

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
		if resp.StatusCode == http.StatusOK {
			currentStatus = statusOperational
		} else {
			currentStatus = statusDisruption
		}
		currentDetail = resp.Status
	}

	st := &status.Status{
		StatusPageId: comp.StatusPageID,
		Component:    comp.ID,
		Container:    comp.ContainerID,
		Details:      currentDetail,
		StatusID:     currentStatus,
	}

	return st, nil
}

func New(timeout, maxRetry int) *Checker {
	return &Checker{maxRetry: maxRetry, timeout: timeout}
}
