package model

import (
	"os"

	"github.com/Jeffail/gabs/v2"
	"gopkg.in/yaml.v2"
)

type Filters []Filter

type Filter struct {
	Testname          string `yaml:"testname"`
	IgnoreJSONPointer string `yaml:"ignoreJSONPointer"`
}

func NewFilterFromFile(filepath string) (*Filters, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	var f Filters
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (f Filters) Filter(testname string, b []byte) ([]byte, error) {
	jsonParsed, err := gabs.ParseJSON(b)
	if err != nil {
		return nil, err
	}
	for _, f := range f {
		if testname != f.Testname {
			continue
		}
		if err := jsonParsed.DeleteP(f.IgnoreJSONPointer); err != nil && err != gabs.ErrNotFound {
			return nil, err
		}
	}
	return jsonParsed.Bytes(), nil
}
