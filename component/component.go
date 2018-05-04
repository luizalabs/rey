package component

import (
	"encoding/json"
	"io/ioutil"
)

type Component struct {
	Name       string `json:"name"`
	HCEndpoint string `json:"hc_endpoint"`
	LastStatus int
	LastDetail string
}

func GetList(path string) ([]*Component, error) {
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	components := make([]*Component, 0)
	if err := json.Unmarshal(jsonFile, &components); err != nil {
		return nil, err
	}
	return components, nil
}
