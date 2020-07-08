//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"fmt"
	"strings"

	sdkmodelapi "github.com/aws/aws-sdk-go/private/model/api"
	oai "github.com/getkin/kin-openapi/openapi3"

	"github.com/jaypipes/aws-api-tools/pkg/model"
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
	sdkAPI    *sdkmodelapi.API
}

func New(serviceAlias string, sdkHelper *model.SDKHelper) (*API, error) {
	modelPath, docPath, err := sdkHelper.ModelAndDocsPath(serviceAlias)
	if err != nil {
		return nil, err
	}
	apiSpec, docSpec, err := parseFrom(modelPath, docPath)
	if err != nil {
		return nil, err
	}
	sdkAPI, err := sdkHelper.API(serviceAlias)
	if err != nil {
		return nil, err
	}
	return &API{
		Alias:    sdkmodelapi.ServiceID(sdkAPI),
		FullName: sdkAPI.Metadata.ServiceFullName,
		Version:  sdkAPI.Metadata.APIVersion,
		Protocol: sdkAPI.Metadata.Protocol,
		apiSpec:  apiSpec,
		docSpec:  docSpec,
		sdkAPI:   sdkAPI,
	}, nil
}

type OperationFilter struct {
	Methods  []string
	Prefixes []string
}

// GetOperations returns a map, keyed by the operation Name/ID, of OpenAPI
// Operation structs
func (a *API) GetOperations(filter *OperationFilter) []*Operation {
	res := []*Operation{}
	filterPrefixes := []string{}
	if filter != nil {
		filterPrefixes = filter.Prefixes
	}
	filterMethods := []string{}
	if filter != nil {
		filterMethods = filter.Methods
	}
	for _, sdkOp := range a.sdkAPI.OperationList() {
		if sdkOp.HTTP.Method == "" {
			continue
		}
		meth := sdkOp.HTTP.Method
		// Match on any of the supplied prefixes
		if len(filterPrefixes) > 0 && !hasAnyPrefix(sdkOp.Name, filterPrefixes) {
			continue
		}
		if len(filterMethods) > 0 && !inStrings(meth, filterMethods) {
			continue
		}
		res = append(res, &Operation{Name: sdkOp.Name, Method: meth})
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
