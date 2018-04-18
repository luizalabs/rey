package aggregator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/luizalabs/rey/status"
)

const (
	statusioURL = "https://api.status.io/v2/component/status/update"
)

type Aggregator struct {
	apiID  string
	apiKey string
}

func (a *Aggregator) Report(s *status.Status) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", statusioURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-id", a.apiID)
	req.Header.Add("x-api-key", a.apiKey)

	cli := new(http.Client)
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error reporting status: %s", string(respBody))
	}

	return nil
}

func New(apiID, apiKey string) *Aggregator {
	return &Aggregator{apiID: apiID, apiKey: apiKey}
}
