//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"strings"

	oai "github.com/getkin/kin-openapi/openapi3"
)

const (
	ObjectTypeObject    = "object"
	ObjectTypeScalar    = "scalar"
	ObjectTypePayload   = "payload"
	ObjectTypeException = "exception"
	ObjectTypeList      = "list"
)

type Object struct {
	Name     string
	Type     string
	DataType string
}

type Operation struct {
	Name   string
	Method string
}

type API struct {
	Alias      string
	FullName   string
	Protocol   string
	Version    string
	apiSpec    *apiSpec
	objectMap  map[string]*Object
	payloads   map[string]bool
	scalars    map[string]bool
	exceptions map[string]bool
	lists      map[string]bool
	swagger    *oai.Swagger
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

// GetOperations returns a map, keyed by the operation Name/ID, of OpenAPI
// Operation structs
func (a *API) GetOperations(filter *OperationFilter) []*Operation {
	a.eval()
	res := []*Operation{}
	filterMethods := []string{}
	if filter != nil {
		filterMethods = filter.Methods
	}
	for _, pathItem := range a.swagger.Paths {
		if pathItem.Get != nil {
			op := pathItem.Get
			// Match on any of the supplied prefixes
			if !hasAnyPrefix(op.OperationID, filter.Prefixes) {
				continue
			}
			if inStrings("GET", filterMethods) {
				res = append(res, &Operation{Name: op.OperationID, Method: "GET"})
				continue
			}
		}
		if pathItem.Head != nil {
			op := pathItem.Get
			// Match on any of the supplied prefixes
			if !hasAnyPrefix(op.OperationID, filter.Prefixes) {
				continue
			}
			if inStrings("HEAD", filterMethods) {
				res = append(res, &Operation{Name: op.OperationID, Method: "HEAD"})
				continue
			}
		}
		if pathItem.Post != nil {
			op := pathItem.Get
			// Match on any of the supplied prefixes
			if !hasAnyPrefix(op.OperationID, filter.Prefixes) {
				continue
			}
			if inStrings("POST", filterMethods) {
				res = append(res, &Operation{Name: op.OperationID, Method: "POST"})
				continue
			}
		}
		if pathItem.Put != nil {
			op := pathItem.Get
			// Match on any of the supplied prefixes
			if !hasAnyPrefix(op.OperationID, filter.Prefixes) {
				continue
			}
			if inStrings("PUT", filterMethods) {
				res = append(res, &Operation{Name: op.OperationID, Method: "PUT"})
				continue
			}
		}
		if pathItem.Delete != nil {
			op := pathItem.Get
			// Match on any of the supplied prefixes
			if !hasAnyPrefix(op.OperationID, filter.Prefixes) {
				continue
			}
			if inStrings("DELETE", filterMethods) {
				res = append(res, &Operation{Name: op.OperationID, Method: "DELETE"})
				continue
			}
		}
		if pathItem.Patch != nil {
			op := pathItem.Get
			// Match on any of the supplied prefixes
			if !hasAnyPrefix(op.OperationID, filter.Prefixes) {
				continue
			}
			if inStrings("PATCH", filterMethods) {
				res = append(res, &Operation{Name: op.OperationID, Method: "PATCH"})
				continue
			}
		}
	}
	return res
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

func (a *API) Schema() *oai.Swagger {
	a.eval()
	info := &oai.Info{
		Title:   a.FullName,
		Version: a.Version,
	}
	exts := map[string]interface{}{}
	info.ExtensionProps = oai.ExtensionProps{exts}
	info.ExtensionProps.Extensions["x-aws-api-alias"] = a.Alias
	info.ExtensionProps.Extensions["x-aws-api-protocol"] = a.Protocol
	a.swagger.Info = info
	a.swagger.OpenAPI = "3.0.0"
	return a.swagger
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
