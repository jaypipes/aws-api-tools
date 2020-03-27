//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import "strings"

type Metadata struct {
	APIVersion      string `json:"apiVersion"`
	ServiceFullName string `json:"serviceFullName"`
}

type InputOutputException struct {
	Shape *string `json:"shape",omitempty`
}

type HTTP struct {
	Method     string `json:"method"`
	RequestURI string `json:"requestUri"`
}

type Operation struct {
	HTTP   HTTP                   `json:"http"`
	Input  InputOutputException   `json:"input"`
	Output InputOutputException   `json:"output"`
	Errors []InputOutputException `json:"errors"`
}

type Shape struct {
	Type      string `json:"type"`
	Exception bool   `json:"exception"`
}

type API struct {
	Metadata   Metadata             `json:"metadata"`
	Operations map[string]Operation `json:"operations"`
	Shapes     map[string]Shape     `json:"shapes"`
}

// GetOperations returns the Shapes in the API that are of a non-compound type
// by returning a map of the shape name and its underlying simple type
func (a API) GetOperations(filterMethod string, prefix string) map[string]Operation {
	res := map[string]Operation{}
	for opName, op := range a.Operations {
		method := op.HTTP.Method
		if filterMethod != "" && filterMethod != method {
			continue
		}
		if prefix != "" && !strings.HasPrefix(opName, prefix) {
			continue
		}
		res[opName] = op
	}
	return res
}

// GetScalars returns the Shapes in the API that are of a non-compound type by
// returning a map of the shape name and its underlying simple type
func (a API) GetScalars() map[string]string {
	res := map[string]string{}
	for shapeName, shape := range a.Shapes {
		if shape.Type != "structure" && shape.Type != "list" {
			res[shapeName] = shape.Type
		}
	}
	return res
}

// GetPayloads returns the Shapes in the API that are used as input or output
// payload wrappers
func (a API) GetPayloads() map[string]Shape {
	res := map[string]Shape{}
	for _, op := range a.Operations {
		if op.Input.Shape != nil {
			inShapeName := *op.Input.Shape
			res[inShapeName] = a.Shapes[inShapeName]
		}
		if op.Output.Shape != nil {
			outShapeName := *op.Output.Shape
			res[outShapeName] = a.Shapes[outShapeName]
		}
	}
	return res
}

// GetObjects returns all Shapes that are structures but are not exception
// classes and that do not look like they are request or response payload
// wrappers.
func (a API) GetObjects() map[string]Shape {
	payloads := a.GetPayloads()
	res := map[string]Shape{}
	for shapeName, shape := range a.Shapes {
		if shape.Type == "structure" && !shape.Exception {
			if _, found := payloads[shapeName]; found {
				continue
			}
			res[shapeName] = shape
		}
	}
	return res
}

// GetExceptions returns all Shapes that are exception classes
func (a API) GetExceptions() map[string]Shape {
	res := map[string]Shape{}
	for shapeName, shape := range a.Shapes {
		if shape.Type == "structure" && shape.Exception {
			res[shapeName] = shape
		}
	}
	return res
}
