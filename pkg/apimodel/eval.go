//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

import (
	oai "github.com/getkin/kin-openapi/openapi3"
)

func (api *API) eval() error {
	if api.swagger != nil {
		return nil
	}

	// Our goal is to convert the CORAL/upstream model into an OpenAPI3
	// specification. We make a number of passes over the raw parsed JSON from
	// the upstream model, attempting to identify what is a scalar, payload,
	// list, etc as well as collating all the operations for the API.
	swagger := newSwagger()
	spec := api.apiSpec
	objectMap := map[string]*Object{}
	comps := &swagger.Components

	for shapeName, shapeSpec := range spec.Shapes {
		var objType string
		// Determine simple types like scalars, lists and exceptions
		if shapeSpec.Type != "structure" && shapeSpec.Type != "list" {
			objType = ObjectTypeScalar
		} else if shapeSpec.Type == "structure" {
			if shapeSpec.Exception {
				objType = ObjectTypeException
			} else {
				// Just a plain ol' object
				objType = ObjectTypeObject
			}
		} else if shapeSpec.Type == "list" {
			objType = ObjectTypeList
		}
		objectMap[shapeName] = &Object{
			Name:     shapeName,
			Type:     objType,
			DataType: shapeSpec.Type,
		}
		schema, err := api.newSchema(shapeName, shapeSpec, []string{})
		if err != nil {
			return err
		}
		comps.Schemas[shapeName] = oai.NewSchemaRef("", schema)
	}
	api.objectMap = objectMap

	for opName, opSpec := range spec.Operations {
		doc := api.docSpec.Operations[opName]
		op, err := opSpec.Operation(opName, doc, swagger, spec)
		if err != nil {
			return err
		}
		// See https://github.com/OAI/OpenAPI-Specification/issues/1635#issuecomment-607444697
		// Many AWS APIs only really use a single HTTP verb (usually POST) and
		// URI (usually /) and vary operations by an "action" parameter.
		// Therefore, for these APIs, we embed a fragment into the URI in order
		// to allow OpenAPI/Swagger to include these as separate operations.
		reqURI := *opSpec.HTTP.RequestURI
		reqURI += "#" + opName
		swagger.AddOperation(reqURI, opSpec.HTTP.Method, op)
		if opSpec.Input != nil && opSpec.Input.ShapeName != nil {
			inShapeName := *opSpec.Input.ShapeName
			shapeSpec := spec.Shapes[inShapeName]
			objectMap[inShapeName] = &Object{
				Name:     inShapeName,
				Type:     ObjectTypePayload,
				DataType: shapeSpec.Type,
			}
		}
		if opSpec.Output != nil && opSpec.Output.ShapeName != nil {
			outShapeName := *opSpec.Output.ShapeName
			shapeSpec := spec.Shapes[outShapeName]
			objectMap[outShapeName] = &Object{
				Name:     outShapeName,
				Type:     ObjectTypePayload,
				DataType: shapeSpec.Type,
			}
		}
	}

	api.swagger = swagger
	return nil
}
