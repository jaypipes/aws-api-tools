//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"strings"
)

type Metadata struct {
	APIVersion      string `json:"apiVersion"`
	ServiceFullName string `json:"serviceFullName"`
}

type ShapeRef struct {
	Shape    *string `json:"shape",omitempty`
	Location *string `json:"location",omitempty`
}

type HTTP struct {
	Method     string `json:"method"`
	RequestURI string `json:"requestUri"`
}

type Operation struct {
	HTTP   HTTP       `json:"http"`
	Input  ShapeRef   `json:"input"`
	Output ShapeRef   `json:"output"`
	Errors []ShapeRef `json:"errors"`
}

type Shape struct {
	Type      string              `json:"type"`
	Exception bool                `json:"exception"`
	Members   map[string]ShapeRef `json:"members"`
}

type API struct {
	Metadata   Metadata             `json:"metadata"`
	Operations map[string]Operation `json:"operations"`
	Shapes     map[string]Shape     `json:"shapes"`
}

// GetOperations returns the Shapes in the API that are of a non-compound type
// by returning a map of the shape name and its underlying simple type
func (a API) GetOperations(filterMethods []string, filterPrefixes []string) map[string]Operation {
	res := map[string]Operation{}
	for opName, op := range a.Operations {
		if len(filterMethods) > 0 {
			// Match on any of the supplied HTTP methods
			found := false
			method := op.HTTP.Method
			for _, filterMethod := range filterMethods {
				if filterMethod == method {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(filterPrefixes) > 0 {
			// Match on any of the supplied prefixes
			found := false
			for _, filterPrefix := range filterPrefixes {
				if strings.HasPrefix(opName, filterPrefix) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
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

func scalarNames(scalars *map[string]string) []string {
	res := make([]string, len(*scalars))
	x := 0
	for scalarName, _ := range *scalars {
		res[x] = scalarName
		x++
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

// GetObjects returns all Shapes that are structures returned from top-level
// operations. Objects are the shapes that are *not* payloads, scalars or
// exceptions
func (a API) GetObjects() map[string]Shape {
	res := map[string]Shape{}
	scalarMap := a.GetScalars()
	scalarNames := scalarNames(&scalarMap)
	createOps := a.GetOperations([]string{}, []string{"Create", "Put"})
	for _, createOp := range createOps {
		if createOp.Output.Shape == nil {
			// Some "create" operations like s3's
			// PutBucketAccelerateConfiguration don't actually create anything
			// but rather modify a specific attribute of an entity and return
			// no content
			continue
		}
		outShapeName := *createOp.Output.Shape
		outShape := a.Shapes[outShapeName]
		for _, shapeRef := range outShape.Members {
			if shapeRef.Shape != nil {
				if shapeRef.Location != nil {
					location := *shapeRef.Location
					if location == "header" || location == "uri" {
						// Some "create" operations like s3's
						// PutObjectRetentionOutput have an output shape whose
						// member shape gets written into the HTTP headers or
						// parts of a URL, not a JSON/XML response element.
						// Ignore these kinds of shapes for the purposes of
						// determining whether the shape is an "object".
						continue
					}
				}
				refShapeName := *shapeRef.Shape
				// Scalars cannot be objects
				if inStrings(refShapeName, scalarNames) {
					continue
				}
				res[refShapeName] = a.Shapes[refShapeName]
			}
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

func inStrings(subject string, collection []string) bool {
	for _, s := range collection {
		if s == subject {
			return true
		}
	}
	return false
}
