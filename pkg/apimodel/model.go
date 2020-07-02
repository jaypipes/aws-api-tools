//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"fmt"
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
	// This is basically the package name in aws-sdk-go/services/
	AliasLower string
	// And this is the sometimes-titlecased alias from the metadata.json file
	// in the aws-sdk-go/services/$alias_lower/$version directory
	Alias     string
	FullName  string
	Protocol  string
	Version   string
	apiSpec   *apiSpec
	docSpec   *docSpec
	objectMap map[string]*Object
	swagger   *oai.Swagger
}

func New(alias string, modelPath string, docPath string) (*API, error) {
	apiSpec, docSpec, err := parseFrom(modelPath, docPath)
	if err != nil {
		return nil, err
	}
	return &API{
		Alias:    apiSpec.Metadata.Alias,
		FullName: apiSpec.Metadata.FullName,
		Version:  apiSpec.Metadata.APIVersion,
		Protocol: apiSpec.Metadata.Protocol,
		apiSpec:  apiSpec,
		docSpec:  docSpec,
	}, nil
}

type OperationFilter struct {
	Methods  []string
	Prefixes []string
}

// GetOperations returns a map, keyed by the operation Name/ID, of OpenAPI
// Operation structs
func (a *API) GetOperations(filter *OperationFilter) []*Operation {
	if err := a.eval(); err != nil {
		fmt.Printf("ERROR evaluating API: %v\n", err)
		return nil
	}
	res := []*Operation{}
	filterPrefixes := []string{}
	if filter != nil {
		filterPrefixes = filter.Prefixes
	}
	filterMethods := []string{}
	if filter != nil {
		filterMethods = filter.Methods
	}
	for _, pathItem := range a.swagger.Paths {
		var op *oai.Operation
		var meth string
		if pathItem.Get != nil {
			op = pathItem.Get
			meth = "GET"
		}
		if pathItem.Head != nil {
			op = pathItem.Head
			meth = "HEAD"
		}
		if pathItem.Post != nil {
			op = pathItem.Post
			meth = "POST"
		}
		if pathItem.Put != nil {
			op = pathItem.Put
			meth = "PUT"
		}
		if pathItem.Delete != nil {
			op = pathItem.Delete
			meth = "DELETE"
		}
		if pathItem.Patch != nil {
			op = pathItem.Patch
			meth = "PATCH"
		}
		// Match on any of the supplied prefixes
		if len(filterPrefixes) > 0 && !hasAnyPrefix(op.OperationID, filterPrefixes) {
			continue
		}
		if len(filterMethods) > 0 && inStrings(meth, filterMethods) {
			continue
		}
		res = append(res, &Operation{Name: op.OperationID, Method: meth})
	}
	return res
}

type ObjectFilter struct {
	Types    []string
	Prefixes []string
}

// GetObjects returns objects that match any of the supplied filter
func (a *API) GetObjects(filter *ObjectFilter) []*Object {
	if err := a.eval(); err != nil {
		fmt.Printf("ERROR evaluating API: %v\n", err)
		return nil
	}
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
	if err := a.eval(); err != nil {
		fmt.Printf("ERROR evaluating API: %v\n", err)
		return nil
	}
	info := &oai.Info{
		Title:       a.FullName,
		Version:     a.Version,
		Description: a.docSpec.Service,
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
