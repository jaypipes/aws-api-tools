//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	"fmt"

	oai "github.com/getkin/kin-openapi/openapi3"
)

func newSwagger() *oai.Swagger {
	swagger := oai.Swagger{}
	comps := oai.NewComponents()
	comps.Schemas = map[string]*oai.SchemaRef{}
	swagger.Components = comps
	return &swagger
}

func (opSpec *opSpec) Operation(opName string, doc string, api *oai.Swagger) (*oai.Operation, error) {
	op := oai.NewOperation()
	op.OperationID = opName
	op.Description = doc

	// Find the shape representing the input to the create operation and
	// add a pointer to an oai.Schema describing the input shape
	if opSpec.Input != nil {
		inShapeName := *opSpec.Input.ShapeName
		_, found := api.Components.Schemas[inShapeName]
		if !found {
			return nil, fmt.Errorf("expected to find input shape schema ref %s", inShapeName)
		}
		inShapeSchemaRef := oai.NewSchemaRef("#/components/schemas/"+inShapeName, nil)
		reqBody := oai.NewRequestBody().WithJSONSchemaRef(inShapeSchemaRef)
		op.RequestBody = &oai.RequestBodyRef{Value: reqBody}
	}
	// Find the shape representing the ouput of the create operation and
	// add fields from the output shape to the resource object, excluding
	// fields already added from the input
	if opSpec.Output != nil {
		outShapeName := *opSpec.Output.ShapeName
		_, found := api.Components.Schemas[outShapeName]
		if !found {
			return nil, fmt.Errorf("expected to find output shape schema ref %s", outShapeName)
		}
		successRespCode := 200
		if opSpec.HTTP.ResponseCode != nil {
			successRespCode = *opSpec.HTTP.ResponseCode
		}
		outShapeSchemaRef := oai.NewSchemaRef("#/components/schemas/"+outShapeName, nil)
		op.AddResponse(successRespCode, oai.NewResponse().WithJSONSchemaRef(outShapeSchemaRef))

		if len(opSpec.Errors) > 0 {
			// TODO(jaypipes): for each error shape, add a response to the
			// Responses object

		}
	}
	return op, nil
}

func newStringSchema(ss *shapeSpec) *oai.Schema {
	schema := oai.NewStringSchema()
	if ss.Min != nil {
		schema.WithMinLength(int64(*ss.Min))
	}
	if ss.Max != nil {
		schema.WithMaxLength(int64(*ss.Max))
	}
	if ss.Pattern != nil {
		schema.WithPattern(*ss.Pattern)
	}
	if len(ss.Enum) > 0 {
		schema.WithEnum(ss.Enum...)
	}
	return schema
}

func newFloat64Schema(ss *shapeSpec) *oai.Schema {
	schema := oai.NewFloat64Schema()
	if ss.Min != nil {
		schema.WithMin(*ss.Min)
	}
	if ss.Max != nil {
		schema.WithMax(*ss.Max)
	}
	return schema
}

func newInt64Schema(ss *shapeSpec) *oai.Schema {
	schema := oai.NewInt64Schema()
	if ss.Min != nil {
		schema.WithMin(*ss.Min)
	}
	if ss.Max != nil {
		schema.WithMax(*ss.Max)
	}
	return schema
}

func (api *API) newArraySchema(
	ss *shapeSpec,
	visitedMemberShapeNames []string,
) (*oai.Schema, error) {
	schema := oai.NewArraySchema()
	if ss.ListMember == nil {
		return nil, fmt.Errorf("expected list member to be non-nil")
	}
	shapeMap := api.apiSpec.Shapes
	listMemberShapeRef := ss.ListMember
	refListMemberShapeName := *listMemberShapeRef.ShapeName
	listMemberShape, found := shapeMap[refListMemberShapeName]
	if !found {
		return nil, fmt.Errorf("expected to find member shape %s", refListMemberShapeName)
	}
	itemsSchema, err := api.newSchema(refListMemberShapeName, listMemberShape, visitedMemberShapeNames)
	if err != nil {
		return nil, err
	}
	schema.WithItems(itemsSchema)
	if ss.Max != nil {
		schema.WithMaxItems(int64(*ss.Max))
	}
	return schema, nil
}

func (api *API) newObjectSchema(
	ss *shapeSpec,
	visitedMemberShapeNames []string,
) (*oai.Schema, error) {
	schema := oai.NewObjectSchema()
	shapeMap := api.apiSpec.Shapes
	for memberName, memberShapeRef := range ss.Members {
		refMemberShapeName := *memberShapeRef.ShapeName
		visitedMemberShapeNames = append(visitedMemberShapeNames, refMemberShapeName)
		memberShape, found := shapeMap[refMemberShapeName]
		if !found {
			return nil, fmt.Errorf("expected to find member shape %s", memberName)
		}
		if memberShape.Type == "structure" && api != nil ||
			inStrings(refMemberShapeName, visitedMemberShapeNames) {
			// If member shape name is in the list of visited member shape
			// names, we detected cycle in member relationships... prevent
			// infinite recursion by just injecting a JSON Reference to the
			// member schema.
			_, refFound := shapeMap[refMemberShapeName]
			if refFound {
				refSchema := oai.NewSchemaRef("#/components/schemas/"+refMemberShapeName, nil)
				schema.WithPropertyRef(memberName, refSchema)
				continue
			}
		}
		memberSchema, err := api.newSchema(refMemberShapeName, memberShape, visitedMemberShapeNames)
		if err != nil {
			return nil, err
		}
		schema.WithProperty(memberName, memberSchema)
	}
	if len(ss.Required) > 0 {
		schema.Required = ss.Required
	}
	return schema, nil
}

// given a shape name, return a new OpenAPI3 Schema representing the shape.
func (api *API) newSchema(
	shapeName string,
	ss *shapeSpec,
	visitedMemberShapeNames []string,
) (*oai.Schema, error) {
	switch ss.Type {
	case "string":
		return newStringSchema(ss), nil
	case "double":
		return newFloat64Schema(ss), nil
	case "long", "integer":
		return newInt64Schema(ss), nil
	case "blob":
		return oai.NewBytesSchema(), nil
	case "boolean":
		return oai.NewBoolSchema(), nil
	case "timestamp":
		return oai.NewDateTimeSchema(), nil
	case "map":
		return oai.NewObjectSchema().WithAnyAdditionalProperties(), nil
	case "list":
		return api.newArraySchema(ss, visitedMemberShapeNames)
	case "structure":
		return api.newObjectSchema(ss, visitedMemberShapeNames)
	}
	return nil, fmt.Errorf("unknown shape type %s", ss.Type)
}
