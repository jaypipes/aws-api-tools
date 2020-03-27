//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

type Metadata struct {
	APIVersion      string `json:"apiVersion"`
	ServiceFullName string `json:"serviceFullName"`
}

type InputOutputException struct {
	Shape *string `json:"shape",omitempty`
}

type Operation struct {
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

// Scalars returns the Shapes in the API that are of a non-compound type by
// returning a map of the shape name and its underlying simple type
func (a API) Scalars() map[string]string {
	res := map[string]string{}
	for shapeName, shape := range a.Shapes {
		if shape.Type != "structure" && shape.Type != "list" {
			res[shapeName] = shape.Type
		}
	}
	return res
}

// Payloads returns the Shapes in the API that are used as input or output
// payload wrappers
func (a API) Payloads() map[string]Shape {
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

// Objects returns all Shapes that are structures but are not exception classes
// and that do not look like they are request or response payload wrappers.
func (a API) Objects() map[string]Shape {
	payloads := a.Payloads()
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

// Exceptions returns all Shapes that are exception classes
func (a API) Exceptions() map[string]Shape {
	res := map[string]Shape{}
	for shapeName, shape := range a.Shapes {
		if shape.Type == "structure" && shape.Exception {
			res[shapeName] = shape
		}
	}
	return res
}
