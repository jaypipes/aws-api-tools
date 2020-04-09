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

func (ss *shapeSpec) Schema(
	shapeName string,
	api *oai.Swagger,
	shapes *map[string]*shapeSpec,
	visitedMemberShapeNames []string,
) *oai.Schema {
	var schema *oai.Schema
	switch ss.Type {
	case "string":
		schema = oai.NewStringSchema()
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
		schema = oai.NewFloat64Schema()
		if ss.Min != nil {
			schema.WithMin(float64(*ss.Min))
		}
		if ss.Max != nil {
			schema.WithMax(float64(*ss.Max))
		}
	case "long", "integer":
		schema = oai.NewInt64Schema()
		if ss.Min != nil {
			schema.WithMin(float64(*ss.Min))
		}
		if ss.Max != nil {
			schema.WithMax(float64(*ss.Max))
		}
	case "blob":
		schema = oai.NewBytesSchema()
	case "boolean":
		schema = oai.NewBoolSchema()
	case "timestamp":
		schema = oai.NewDateTimeSchema()
	case "map":
		schema = oai.NewObjectSchema().WithAnyAdditionalProperties()
	case "list":
		schema = oai.NewArraySchema()
		if shapes == nil {
			panic("expected shapes to be non-nil")
		}
		if ss.ListMember == nil {
			panic("expected list member to be non-nil")
		}
		shapeMap := *shapes
		listMemberShapeRef := ss.ListMember
		refListMemberShapeName := *listMemberShapeRef.ShapeName
		listMemberShape, found := shapeMap[refListMemberShapeName]
		if !found {
			panic("expected to find member shape " + refListMemberShapeName)
		}
		itemsSchema := listMemberShape.Schema(refListMemberShapeName, api, shapes, visitedMemberShapeNames)
		schema.WithItems(itemsSchema)
		if ss.Max != nil {
			schema.WithMaxItems(*ss.Max)
		}
	case "structure":
		schema = oai.NewObjectSchema()
		if shapes == nil {
			panic("expected shapes to be non-nil")
		}
		shapeMap := *shapes
		for memberName, memberShapeRef := range ss.Members {
			refMemberShapeName := *memberShapeRef.ShapeName
			visitedMemberShapeNames = append(visitedMemberShapeNames, refMemberShapeName)
			memberShape, found := shapeMap[refMemberShapeName]
			if !found {
				panic("expected to find member shape " + memberName)
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
			schema.WithProperty(memberName, memberShape.Schema(refMemberShapeName, api, shapes, visitedMemberShapeNames))
		}
		if len(ss.Required) > 0 {
			schema.Required = ss.Required
		}
	}
	return schema
}
