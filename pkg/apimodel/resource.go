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

// Many service APIs follow a pattern that we can use to determine top-level or
// resource objects:
//
// There will be a Create operation that involves the resource object called
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
// We will be able to identify the fields in the resource object by looking at
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

func getResources(api *API) (map[string]*Resource, error) {
	pluralize := pluralize.NewClient()
	apiProtocol := api.Metadata.Protocol
	resources := map[string]*Resource{}
	filter := &OperationFilter{
		Prefixes: []string{"Create"},
	}
	createOps := api.GetOperations(filter)
	for _, createOp := range createOps {
		// Some API operations are "CreateOrUpdate" -- i.e. "Replace". For
		// instance, the AWS Autoscaling API has a CreateOrUpdateTags
		// operation. Trim off the "CreateOrUpdate" prefix for these
		// operations.
		var objName string
		if strings.HasPrefix(createOp.Name, "CreateOrUpdate") {
			objName = strings.TrimPrefix(createOp.Name, "CreateOrUpdate")
		} else {
			objName = strings.TrimPrefix(createOp.Name, "Create")
		}
		singularName := pluralize.Singular(objName)

		// Tag is a special case. It is often represented as a
		// top-level/resource object because there are CreateOrUpdateTags
		// operations that accept a payload that replaces all tags on a
		// specific resource. However, Tag is not an actual resource object.
		// Instead, nearly all resources can have zero or more key/value pairs
		// associated with them (these are tags).
		if singularName == "Tag" {
			continue
		}

		// For APIs that have a "rest-json" protocol, we can look at the operation's
		// http.requestUri field to determine whether the operation is on a "top-level"
		// object.
		if apiProtocol == "rest-json" {
			if strings.Count(createOp.RequestURI, "/") > 1 {
				continue
			}
		}

		pluralName := pluralize.Plural(objName)
		resource := &Resource{
			SingularName: singularName,
			PluralName:   pluralName,
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
				resource.Fields[fieldName] = &Field{
					Type:       field.Type,
					IsRequired: inStrings(fieldName, field.RequiredFieldNames),
					IsMutable:  true,
				}
			}
		}
		// Find the shape representing the ouput of the create operation and
		// add fields from the output shape to the resource object, excluding
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
				if _, found := resource.Fields[fieldName]; !found {
					resource.Fields[fieldName] = &Field{
						Type:       field.Type,
						IsRequired: inStrings(field.Name, field.RequiredFieldNames),
						IsMutable:  false,
					}
				}
			}
		}
		resources[objName] = resource
	}
	return resources, nil
}
