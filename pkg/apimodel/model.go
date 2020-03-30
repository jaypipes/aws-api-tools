//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"strings"
)

type Shape struct {
	Name               string
	Type               string
	Fields             map[string]*Shape
	RequiredFieldNames []string
}

type Operation struct {
	Name   string
	Method string
	Input  *Shape
	Output *Shape
	Errors []*Shape
}

type Field struct {
	Type       string
	IsRequired bool
	IsMutable  bool
}

type Resource struct {
	SingularName string
	PluralName   string
	Fields       map[string]*Field
}

type API struct {
	Metadata     metadataSpec
	opMap        map[string]Operation
	shapeMap     map[string]Shape
	payloadMap   map[string]*Shape
	scalarMap    map[string]*Shape
	exceptionMap map[string]*Shape
	objectMap    map[string]*Shape
	listMap      map[string]*Shape
	resourceMap  map[string]*Resource
}

type OperationFilter struct {
	Methods  []string
	Prefixes []string
}

// GetOperations returns the Shapes in the API that are of a non-compound type
// by returning a map of the shape name and its underlying simple type
func (a *API) GetOperations(filter *OperationFilter) []*Operation {
	res := []*Operation{}
	for opName, op := range a.opMap {
		if filter != nil {
			if len(filter.Methods) > 0 {
				// Match on any of the supplied HTTP methods
				if !inStrings(op.Method, filter.Methods) {
					continue
				}
			}
			if len(filter.Prefixes) > 0 {
				// Match on any of the supplied prefixes
				if !hasAnyPrefix(opName, filter.Prefixes) {
					continue
				}
			}
		}
		resOp := a.opMap[opName]
		res = append(res, &resOp)
	}
	return res
}

// GetScalars returns the Shapes in the API that are of a non-compound type by
// returning a map of the shape name and its underlying simple type
func (a *API) GetScalars() []*Shape {
	res := make([]*Shape, len(a.scalarMap))
	x := 0
	for _, scalar := range a.scalarMap {
		res[x] = scalar
		x++
	}
	return res
}

// GetPayloads returns the Shapes in the API that are used as input or output
// payload wrappers
func (a *API) GetPayloads() []*Shape {
	res := make([]*Shape, len(a.payloadMap))
	x := 0
	for _, payload := range a.payloadMap {
		res[x] = payload
		x++
	}
	return res
}

// GetResources returns objects that have been identified as top-level resource
// structures for the API.
func (a *API) GetResources() []*Resource {
	res := make([]*Resource, len(a.resourceMap))
	x := 0
	for _, resource := range a.resourceMap {
		res[x] = resource
		x++
	}
	return res
}

// GetObjects returns shapes that are *not* payloads, scalars, lists or
// exceptions
func (a *API) GetObjects() []*Shape {
	res := make([]*Shape, len(a.objectMap))
	x := 0
	for _, object := range a.objectMap {
		res[x] = object
		x++
	}
	return res
}

// GetExceptions returns all Shapes that are exception classes
func (a *API) GetExceptions() []*Shape {
	res := make([]*Shape, len(a.exceptionMap))
	x := 0
	for _, exception := range a.exceptionMap {
		res[x] = exception
		x++
	}
	return res
}

// GetLists returns all Shapes that are list classes
func (a *API) GetLists() []*Shape {
	res := make([]*Shape, len(a.listMap))
	x := 0
	for _, list := range a.listMap {
		res[x] = list
		x++
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

func hasAnyPrefix(subject string, prefixes []string) bool {
	// Match on any of the supplied prefixes
	for _, prefix := range prefixes {
		if strings.HasPrefix(subject, prefix) {
			return true
		}
	}
	return false
}
