//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type metadataSpec struct {
	APIVersion string `json:"apiVersion"`
	FullName   string `json:"serviceFullName"`
	Protocol   string `json:"protocol"`
}

type shapeRefSpec struct {
	ShapeName *string `json:"shape,omitempty"`
	Location  *string `json:"location,omitempty"`
}

type httpSpec struct {
	Method       string  `json:"method"`
	RequestURI   *string `json:"requestUri"`
	ResponseCode *int    `json:"responseCode"`
}

type opSpec struct {
	HTTP   *httpSpec       `json:"http,omitempty"`
	Input  *shapeRefSpec   `json:"input,omitempty"`
	Output *shapeRefSpec   `json:"output,omitempty"`
	Errors []*shapeRefSpec `json:"errors"`
}

type shapeSpec struct {
	Type       string                   `json:"type"`
	Exception  bool                     `json:"exception"`
	Required   []string                 `json:"required"`
	Members    map[string]*shapeRefSpec `json:"members"`
	ListMember *shapeRefSpec            `json:"member,omitempty"` // for list types
	Min        *float64                 `json:"min,omitempty"`
	Max        *float64                 `json:"max,omitempty"`
	Pattern    *string                  `json:"pattern,omitempty"`
	Enum       []interface{}            `json:"enum"`
}

type apiSpec struct {
	Metadata   metadataSpec          `json:"metadata"`
	Operations map[string]*opSpec    `json:"operations"`
	Shapes     map[string]*shapeSpec `json:"shapes"`
}

type shapeDocSpec struct {
	Base *string           `json:"base"`
	Refs map[string]string `json:"refs"`
}

type docSpec struct {
	Service    string                   `json:"service"`
	Operations map[string]string        `json:"operations"`
	Shapes     map[string]*shapeDocSpec `json:"shapes"`
}

func parseFrom(modelPath string, docPath string) (*apiSpec, *docSpec, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("expected to find %s", modelPath)
	}
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("expected to find %s", docPath)
	}
	var apiSpec apiSpec
	b, err := ioutil.ReadFile(modelPath)
	if err = json.Unmarshal(b, &apiSpec); err != nil {
		return nil, nil, err
	}
	var docSpec docSpec
	b, err = ioutil.ReadFile(docPath)
	if err = json.Unmarshal(b, &docSpec); err != nil {
		return nil, nil, err
	}
	return &apiSpec, &docSpec, nil
}
