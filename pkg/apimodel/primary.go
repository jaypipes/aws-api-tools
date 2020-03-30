//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"fmt"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
)

// If a service API has a protocol of "query", the API generally follows a
// pattern that we can use to determine top-level or primary objects:
//
// There will be a Create operation that involves the primary object called
// Create{$ObjectName}. An example of this from the SNS API:
//
//  "CreateTopic":{
//    "name":"CreateTopic",
//    "http":{
//      "method":"POST",
//      "requestUri":"/"
//    },
//    "input":{"shape":"CreateTopicInput"},
//    "output":{
//      "shape":"CreateTopicResponse",
//      "resultWrapper":"CreateTopicResult"
//    },
//    "errors":[
//      {"shape":"InvalidParameterException"},
//      {"shape":"TopicLimitExceededException"},
//      {"shape":"InternalErrorException"},
//      {"shape":"AuthorizationErrorException"},
//      {"shape":"InvalidSecurityException"},
//      {"shape":"TagLimitExceededException"},
//      {"shape":"StaleTagException"},
//      {"shape":"TagPolicyException"},
//      {"shape":"ConcurrentAccessException"}
//    ]
//  },
//
// We will be able to identify the fields in the primary object by looking at
// the input Shape and grabbing the shape that is listed in its single member.
// For example, the CreateTopicInput shape from the SNS API:
//
//  "CreateTopicInput":{
//    "type":"structure",
//    "required":["Name"],
//    "members":{
//      "Name":{"shape":"topicName"},
//      "Attributes":{"shape":"TopicAttributesMap"},
//      "Tags":{"shape":"TagList"}
//    }
//  },
//
// In addition, if we examine the output Shape from the Create operation, we
// will typically be able to determine how the object is expected to be
// identified. In the case of the SNS API, it is via an ARN:
//
//  "CreateTopicResponse":{
//    "type":"structure",
//    "members":{
//      "TopicArn":{"shape":"topicARN"}
//    }
//  },

func getQueryProtocolPrimaries(api *API) (map[string]*Primary, error) {
	pluralize := pluralize.NewClient()
	primaries := map[string]*Primary{}
	filter := &OperationFilter{
		Prefixes: []string{"Create"},
	}
	createOps := api.GetOperations(filter)
	for _, createOp := range createOps {
		objName := strings.TrimPrefix(createOp.Name, "Create")
		primary := &Primary{
			SingularName: objName,
			PluralName:   pluralize.Plural(objName),
			Fields:       map[string]*Field{},
		}
		// Find the shape representing the input to the create operation and
		// add fields from the input shape
		if createOp.Input != nil {
			inShapeName := createOp.Input.Name
			inShape, found := api.shapeMap[inShapeName]
			if !found {
				return nil, fmt.Errorf("expected to find input shape %s", inShapeName)
			}
			if inShape.Type != "structure" {
				return nil, fmt.Errorf("expected to find a structure type for input shape %s but found %s", inShapeName, inShape.Type)
			}
			for fieldName, field := range inShape.Fields {
				primary.Fields[fieldName] = &Field{
					Type:       field.Type,
					IsRequired: inStrings(fieldName, field.RequiredFieldNames),
					IsMutable:  true,
				}
			}
		}
		// Find the shape representing the ouput of the create operation and
		// add fields from the output shape to the primary object, excluding
		// fields already added from the input
		if createOp.Output != nil {
			outShapeName := createOp.Output.Name
			outShape, found := api.shapeMap[outShapeName]
			if !found {
				return nil, fmt.Errorf("expected to find input shape %s", outShapeName)
			}
			if outShape.Type != "structure" {
				return nil, fmt.Errorf("expected to find a structure type for input shape %s but found %s", outShapeName, outShape.Type)
			}
			for fieldName, field := range outShape.Fields {
				if _, found := primary.Fields[fieldName]; !found {
					primary.Fields[fieldName] = &Field{
						Type:       field.Type,
						IsRequired: inStrings(field.Name, field.RequiredFieldNames),
						IsMutable:  false,
					}
				}
			}
		}
		primaries[objName] = primary
	}
	return primaries, nil
}
