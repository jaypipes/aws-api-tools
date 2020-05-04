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
	// Helpfully decorate the object with an annotation if this object is an
	// exception type
	if ss.Exception {
		exts := map[string]interface{}{}
		schema.ExtensionProps = oai.ExtensionProps{exts}
		schema.ExtensionProps.Extensions["x-aws-api-exception"] = true
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
	case "double", "float":
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
