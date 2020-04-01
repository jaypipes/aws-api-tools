//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	oai "github.com/getkin/kin-openapi/openapi3"
)

type Resource struct {
	SingularName string
	PluralName   string
	createOp     *oai.Operation
}

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

//func getResources(api *API) (map[string]*Resource, error) {
//	pluralize := pluralize.NewClient()
//	resources := map[string]*Resource{}
//	filter := &OperationFilter{
//		Prefixes: []string{"Create"},
//	}
//	createOps := api.GetOperations(filter)
//	for _, createOp := range createOps {
//		// Some API operations are "CreateOrUpdate" -- i.e. "Replace". For
//		// instance, the AWS Autoscaling API has a CreateOrUpdateTags
//		// operation. Trim off the "CreateOrUpdate" prefix for these
//		// operations.
//		var objName string
//		if strings.HasPrefix(createOp.Name, "CreateOrUpdate") {
//			objName = strings.TrimPrefix(createOp.Name, "CreateOrUpdate")
//		} else {
//			objName = strings.TrimPrefix(createOp.Name, "Create")
//		}
//		singularName := pluralize.Singular(objName)
//
//		// Tag is a special case. It is often represented as a
//		// top-level/resource object because there are CreateOrUpdateTags
//		// operations that accept a payload that replaces all tags on a
//		// specific resource. However, Tag is not an actual resource object.
//		// Instead, nearly all resources can have zero or more key/value pairs
//		// associated with them (these are tags).
//		if singularName == "Tag" {
//			continue
//		}
//
//		pluralName := pluralize.Plural(objName)
//		createOAIOp := oai.NewOperation()
//		resource := &Resource{
//			SingularName: singularName,
//			PluralName:   pluralName,
//		}
//		// Find the shape representing the input to the create operation and
//		// add a pointer to an oai.Schema describing the input shape
//		if createOp.Input != nil {
//			inShapeName := createOp.Input.Name
//			inObj, found := api.objectMap[inShapeName]
//			if !found {
//				return nil, fmt.Errorf("expected to find input shape %s", inShapeName)
//			}
//			if inObj.DataType != "structure" {
//				return nil, fmt.Errorf("expected to find a structure type for input shape %s but found %s", inShapeName, inObj.Type)
//			}
//			reqBody := oai.NewRequestBody().WithJSONSchema(inObj.Schema())
//			createOAIOp.OperationID = createOp.Name
//			createOAIOp.RequestBody = &oai.RequestBodyRef{Value: reqBody}
//			api.swagger.AddOperation(createOp.opSpec.HTTP.RequestURI, createOp.opSpec.HTTP.Method, createOAIOp)
//		}
//		// Find the shape representing the ouput of the create operation and
//		// add fields from the output shape to the resource object, excluding
//		// fields already added from the input
//		if createOp.Output != nil {
//			outShapeName := createOp.Output.Name
//			outObj, found := api.objectMap[outShapeName]
//			if !found {
//				return nil, fmt.Errorf("expected to find output shape %s", outShapeName)
//			}
//			if outObj.DataType != "structure" {
//				return nil, fmt.Errorf("expected to find a structure type for output shape %s but found %s", outShapeName, outObj.Type)
//			}
//			successRespCode := 200
//			if createOp.opSpec.HTTP.ResponseCode != nil {
//				successRespCode = *createOp.opSpec.HTTP.ResponseCode
//			}
//			createOAIOp.AddResponse(successRespCode, oai.NewResponse().WithJSONSchema(outObj.Schema()))
//
//			if len(createOp.Errors) > 0 {
//				// TODO(jaypipes): for each error shape, add a response to the
//				// Responses object
//				continue
//			}
//		}
//		resource.createOp = createOAIOp
//		resources[objName] = resource
//	}
//	return resources, nil
//}
