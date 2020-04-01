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
	"github.com/getkin/kin-openapi/openapi3"
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
	apiProtocol := api.Protocol
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
			Properties:   map[string]*openapi3.Schema{},
		}
		// Find the shape representing the input to the create operation and
		// add a pointer to an openapi3.Schema describing the input shape
		if createOp.Input != nil {
			inShapeName := createOp.Input.Name
			inObj, found := api.objectMap[inShapeName]
			if !found {
				return nil, fmt.Errorf("expected to find input shape %s", inShapeName)
			}
			if inObj.DataType != "structure" {
				return nil, fmt.Errorf("expected to find a structure type for input shape %s but found %s", inShapeName, inObj.Type)
			}
			for memberName, member := range inObj.Members {
				resource.Properties[memberName] = member.OpenAPISchema()
			}
			// This is about as good as we're gonna get w.r.t. determining
			// which fields of a resource are actually required... look at the
			// Input shape's required fields, if any.
			resource.Required = inObj.shapeSpec.Required
		}
		// Find the shape representing the ouput of the create operation and
		// add fields from the output shape to the resource object, excluding
		// fields already added from the input
		if createOp.Output != nil {
			outShapeName := createOp.Output.Name
			outObj, found := api.objectMap[outShapeName]
			if !found {
				return nil, fmt.Errorf("expected to find output shape %s", outShapeName)
			}
			if outObj.DataType != "structure" {
				return nil, fmt.Errorf("expected to find a structure type for output shape %s but found %s", outShapeName, outObj.Type)
			}

			var membersToProcess *map[string]*Object = &outObj.Members

			// Often (but annoyingly not always), the API's response will wrap
			// the returned members with a single wrapper member named the same
			// as the resource. For example, the EKS CreateCluster operation's
			// output Shape looks like this:
			//
			// "CreateClusterResponse":{
			//   "type":"structure",
			//   "members":{
			//     "cluster":{"shape":"Cluster"}
			//   }
			// },
			//
			// If this is the case, we go ahead and "flatten" things by
			// processing the single member's members...
			if len(outObj.Members) == 1 {
				for memberShapeName, memberObj := range outObj.Members {
					if strings.ToLower(singularName) == strings.ToLower(memberShapeName) {
						membersToProcess = &memberObj.Members
						break
					}
				}
			}

			for memberName, member := range *membersToProcess {
				if _, found := resource.Properties[memberName]; !found {
					resource.Properties[memberName] = member.OpenAPISchema()
				}
			}
		}
		resources[objName] = resource
	}
	return resources, nil
}

func (object *Object) OpenAPISchema() *openapi3.Schema {
	var schema *openapi3.Schema
	ss := object.shapeSpec
	switch object.DataType {
	case "string":
		schema = openapi3.NewStringSchema()
		if ss.Min != nil {
			schema.WithMinLength(*ss.Min)
		}
		if ss.Max != nil {
			schema.WithMaxLength(*ss.Max)
		}
		if ss.Pattern != nil {
			schema.WithPattern(*ss.Pattern)
		}
		if len(ss.Enum) > 0 {
			schema.WithEnum(ss.Enum...)
		}
	case "double":
		schema = openapi3.NewFloat64Schema()
		if ss.Min != nil {
			schema.WithMin(float64(*ss.Min))
		}
		if ss.Max != nil {
			schema.WithMax(float64(*ss.Max))
		}
	case "long", "integer":
		schema = openapi3.NewInt64Schema()
		if ss.Min != nil {
			schema.WithMin(float64(*ss.Min))
		}
		if ss.Max != nil {
			schema.WithMax(float64(*ss.Max))
		}
	case "blob":
		schema = openapi3.NewBytesSchema()
	case "boolean":
		schema = openapi3.NewBoolSchema()
	case "timestamp":
		schema = openapi3.NewDateTimeSchema()
	case "map":
		schema = openapi3.NewObjectSchema().WithAnyAdditionalProperties()
	case "list":
		schema = openapi3.NewArraySchema()
		for _, memberObj := range object.Members {
			itemsSchema := memberObj.OpenAPISchema()
			schema.WithItems(itemsSchema)
			break
		}
		if ss.Max != nil {
			schema.WithMaxItems(*ss.Max)
		}
	case "structure":
		schema = openapi3.NewObjectSchema()
		memberProps := map[string]*openapi3.Schema{}
		for memberName, memberObj := range object.Members {
			memberProps[memberName] = memberObj.OpenAPISchema()
		}
		schema.WithProperties(memberProps)
		if len(ss.Required) > 0 {
			schema.Required = ss.Required
		}
	}
	return schema
}

func (r *Resource) OpenAPI3Schema() *openapi3.Schema {
	schema := openapi3.NewObjectSchema().WithProperties(r.Properties)
	if len(r.Required) > 0 {
		schema.Required = r.Required
	}
	return schema
}
