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

	oai "github.com/getkin/kin-openapi/openapi3"
)

type metadataSpec struct {
	APIVersion string `json:"apiVersion"`
	FullName   string `json:"serviceFullName"`
	Protocol   string `json:"protocol"`
}

type shapeRefSpec struct {
	ShapeName *string `json:"shape",omitempty`
	Location  *string `json:"location",omitempty`
}

type httpSpec struct {
	Method       string  `json:"method"`
	RequestURI   *string `json:"requestUri"`
	ResponseCode *int    `json:"responseCode"`
}

type opSpec struct {
	HTTP   *httpSpec       `json:"http",omitempty`
	Input  *shapeRefSpec   `json:"input",omitempty`
	Output *shapeRefSpec   `json:"output",omitempty`
	Errors []*shapeRefSpec `json:"errors"`
}

type shapeSpec struct {
	Type       string                   `json:"type"`
	Exception  bool                     `json:"exception"`
	Required   []string                 `json:"required"`
	Members    map[string]*shapeRefSpec `json:"members"`
	ListMember *shapeRefSpec            `json:"member",omitempty` // for list types
	Min        *int64                   `json:"min",omitempty`
	Max        *int64                   `json:"max",omitempty`
	Pattern    *string                  `json:"pattern",omitempty`
	Enum       []interface{}            `json:"enum"`
}

type apiSpec struct {
	Metadata   metadataSpec          `json:"metadata"`
	Operations map[string]*opSpec    `json:"operations"`
	Shapes     map[string]*shapeSpec `json:"shapes"`
}

func parseFrom(modelPath string) (*apiSpec, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("expected to find %s", modelPath)
	}
	var spec apiSpec
	b, err := ioutil.ReadFile(modelPath)
	if err = json.Unmarshal(b, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func (api *API) eval() error {
	if api.swagger != nil {
		return nil
	}
	api.objectMap = map[string]*Object{}

	// Our goal is to convert the CORAL/upstream model into an OpenAPI3
	// specification. We make a number of passes over the raw parsed JSON from
	// the upstream model, attempting to identify what is a scalar, payload,
	// list, etc as well as collating all the operations for the API.
	swagger := &oai.Swagger{}
	comps := oai.NewComponents()
	comps.Schemas = map[string]*oai.SchemaRef{}
	swagger.Components = comps
	spec := api.apiSpec
	api.payloads = map[string]bool{}
	api.scalars = map[string]bool{}
	api.exceptions = map[string]bool{}
	api.lists = map[string]bool{}

	// Populate the base object maps
	for shapeName, shapeSpec := range spec.Shapes {
		comps.Schemas[shapeName] = &oai.SchemaRef{Value: shapeSpec.Schema(swagger, &spec.Shapes)}
		// Determine simple types like scalars, lists and exceptions
		if shapeSpec.Type != "structure" && shapeSpec.Type != "list" {
			api.scalars[shapeName] = true
			api.objectMap[shapeName] = &Object{
				Name:     shapeName,
				Type:     ObjectTypeScalar,
				DataType: shapeSpec.Type,
			}
		} else if shapeSpec.Type == "structure" && shapeSpec.Exception {
			api.exceptions[shapeName] = true
			api.objectMap[shapeName] = &Object{
				Name:     shapeName,
				Type:     ObjectTypeException,
				DataType: shapeSpec.Type,
			}
		} else if shapeSpec.Type == "list" {
			api.lists[shapeName] = true
			api.objectMap[shapeName] = &Object{
				Name:     shapeName,
				Type:     ObjectTypeList,
				DataType: shapeSpec.Type,
			}
		}
	}

	for opName, opSpec := range spec.Operations {
		op, err := opSpec.Operation(opName, swagger)
		if err != nil {
			return err
		}
		// See https://github.com/OAI/OpenAPI-Specification/issues/1635#issuecomment-607444697
		// Some AWS APIs (those that are of the "query" protocol, only really
		// use a single HTTP verb and URI (usually /) and vary operations by
		// the "action" parameter. Therefore, for these APIs, we embed a
		// fragment into the URI in order to allow OpenAPI/Swagger to include
		// these as separate operations.
		reqURI := *opSpec.HTTP.RequestURI
		if api.Protocol == "query" {
			reqURI += "#action=" + opName
		}
		swagger.AddOperation(reqURI, opSpec.HTTP.Method, op)
		if opSpec.Input != nil && opSpec.Input.ShapeName != nil {
			inShapeName := *opSpec.Input.ShapeName
			shapeSpec := spec.Shapes[inShapeName]
			api.payloads[inShapeName] = true
			api.objectMap[inShapeName] = &Object{
				Name:     inShapeName,
				Type:     ObjectTypePayload,
				DataType: shapeSpec.Type,
			}
		}
		if opSpec.Output != nil && opSpec.Output.ShapeName != nil {
			outShapeName := *opSpec.Output.ShapeName
			shapeSpec := spec.Shapes[outShapeName]
			api.payloads[outShapeName] = true
			api.objectMap[outShapeName] = &Object{
				Name:     outShapeName,
				Type:     ObjectTypePayload,
				DataType: shapeSpec.Type,
			}
		}
	}

	api.swagger = swagger
	return nil
}
