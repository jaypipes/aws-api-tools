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
	APIVersion      string `json:"apiVersion"`
	ServiceFullName string `json:"serviceFullName"`
}

type shapeRefSpec struct {
	ShapeName *string `json:"shape",omitempty`
	Location  *string `json:"location",omitempty`
}

type httpSpec struct {
	Method     string `json:"method"`
	RequestURI string `json:"requestUri"`
}

type opSpec struct {
	HTTP   httpSpec       `json:"http"`
	Input  shapeRefSpec   `json:"input"`
	Output shapeRefSpec   `json:"output"`
	Errors []shapeRefSpec `json:"errors"`
}

type shapeSpec struct {
	Type      string                  `json:"type"`
	Exception bool                    `json:"exception"`
	Members   map[string]shapeRefSpec `json:"members"`
}

type apiSpec struct {
	Metadata   metadataSpec         `json:"metadata"`
	Operations map[string]opSpec    `json:"operations"`
	Shapes     map[string]shapeSpec `json:"shapes"`
}

func ParseFrom(modelPath string) (*API, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("expected to find %s", modelPath)
	}
	var spec apiSpec
	b, err := ioutil.ReadFile(modelPath)
	if err = json.Unmarshal(b, &spec); err != nil {
		return nil, err
	}

	api, err := apiFromSpec(&spec)
	if err != nil {
		return nil, err
	}
	return api, nil
}

func apiFromSpec(spec *apiSpec) (*API, error) {
	api := API{
		Metadata:     spec.Metadata,
		shapeMap:     make(map[string]Shape, len(spec.Shapes)),
		opMap:        make(map[string]Operation, len(spec.Operations)),
		payloadMap:   map[string]*Shape{},
		scalarMap:    map[string]*Shape{},
		exceptionMap: map[string]*Shape{},
		objectMap:    map[string]*Shape{},
		listMap:      map[string]*Shape{},
	}
	// Populate the base shape and operation maps
	for shapeName, shapeSpec := range spec.Shapes {
		shape := Shape{
			Name: shapeName,
			Type: shapeSpec.Type,
		}
		api.shapeMap[shapeName] = shape

		// Determine simple types like scalars, lists and exceptions
		if shapeSpec.Type != "structure" && shapeSpec.Type != "list" {
			api.scalarMap[shapeName] = &shape
		} else if shapeSpec.Type == "structure" && shapeSpec.Exception {
			api.exceptionMap[shapeName] = &shape
		} else if shapeSpec.Type == "list" {
			api.listMap[shapeName] = &shape
		}
	}
	for opName, opSpec := range spec.Operations {
		api.opMap[opName] = Operation{
			Name:   opName,
			Method: opSpec.HTTP.Method,
		}
		// Determine payload types by examining the Input and Output pointers
		// for Operations
		for _, opSpec := range spec.Operations {
			if opSpec.Input.ShapeName != nil {
				inShapeName := *opSpec.Input.ShapeName
				sh, ok := api.shapeMap[inShapeName]
				if !ok {
					return nil, fmt.Errorf("expected to find shape %s", inShapeName)
				}
				api.payloadMap[inShapeName] = &sh

			}
			if opSpec.Output.ShapeName != nil {
				outShapeName := *opSpec.Output.ShapeName
				sh, ok := api.shapeMap[outShapeName]
				if !ok {
					return nil, fmt.Errorf("expected to find shape %s", outShapeName)
				}
				api.payloadMap[outShapeName] = &sh
			}
		}
	}

	// objects are the shapes that are *not* payloads, scalars or exceptions
	for shapeName, _ := range spec.Shapes {
		if _, found := api.scalarMap[shapeName]; found {
			continue
		}
		if _, found := api.payloadMap[shapeName]; found {
			continue
		}
		if _, found := api.listMap[shapeName]; found {
			continue
		}
		if _, found := api.exceptionMap[shapeName]; found {
			continue
		}
		objShape := api.shapeMap[shapeName]
		api.objectMap[shapeName] = &objShape
	}
	return &api, nil
}
