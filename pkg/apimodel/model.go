//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	ObjectTypeObject    = "object"
	ObjectTypeScalar    = "scalar"
	ObjectTypePayload   = "payload"
	ObjectTypeException = "exception"
	ObjectTypeList      = "list"
)

type Object struct {
	Name      string
	Type      string
	DataType  string
	Members   map[string]*Object
	shapeSpec shapeSpec
}

type Operation struct {
	Name       string
	Method     string
	RequestURI string
	Input      *Object
	Output     *Object
	Errors     []*Object
}

type Resource struct {
	SingularName string
	PluralName   string
	Properties   map[string]*openapi3.Schema
	Required     []string
}

type API struct {
	Alias        string
	FullName     string
	Protocol     string
	Version      string
	apiSpec      *apiSpec
	opMap        map[string]Operation
	objectMap    map[string]*Object
	payloadMap   map[string]*Object
	scalarMap    map[string]*Object
	exceptionMap map[string]*Object
	listMap      map[string]*Object
	resourceMap  map[string]*Resource
}

func New(alias string, modelPath string) (*API, error) {
	spec, err := parseFrom(modelPath)
	if err != nil {
		return nil, err
	}
	return &API{
		Alias:    alias,
		FullName: spec.Metadata.FullName,
		Version:  spec.Metadata.APIVersion,
		Protocol: spec.Metadata.Protocol,
		apiSpec:  spec,
	}, nil
}

type OperationFilter struct {
	Methods  []string
	Prefixes []string
}

// GetOperations returns the Shapes in the API that are of a non-compound type
// by returning a map of the shape name and its underlying simple type
func (a *API) GetOperations(filter *OperationFilter) []*Operation {
	a.eval()
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

// GetResources returns objects that have been identified as top-level resource
// structures for the API.
func (a *API) GetResources() []*Resource {
	a.eval()
	res := make([]*Resource, len(a.resourceMap))
	x := 0
	for _, resource := range a.resourceMap {
		res[x] = resource
		x++
	}
	return res
}

// GetResource returns a Resource with the specified name
func (a *API) GetResource(resourceName string) (*Resource, error) {
	a.eval()
	r, found := a.resourceMap[resourceName]
	if !found {
		return nil, fmt.Errorf("no such resource '%s'", resourceName)
	}
	return r, nil
}

type ObjectFilter struct {
	Types    []string
	Prefixes []string
}

// GetObjects returns objects that match any of the supplied filter
func (a *API) GetObjects(filter *ObjectFilter) []*Object {
	a.eval()
	res := []*Object{}
	for objectName, object := range a.objectMap {
		if filter != nil {
			if len(filter.Types) > 0 {
				// Match on any of the supplied object types
				if !inStrings(object.Type, filter.Types) {
					continue
				}
			}
			if len(filter.Prefixes) > 0 {
				// Match on any of the supplied prefixes
				if !hasAnyPrefix(objectName, filter.Prefixes) {
					continue
				}
			}
		}
		res = append(res, object)
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
