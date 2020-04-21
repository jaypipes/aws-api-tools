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

func (opSpec *opSpec) Operation(opName string, doc string, api *oai.Swagger, apiSpec *apiSpec) (*oai.Operation, error) {
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
			// Due to the way OpenAPI3 Operations are structured, we need to
			// first gather all errors that occur with the same HTTP response
			// code and then use an oneOf JSONSchema to inject each of the
			// possible error responses into the Operation
			codeSchemaMap := map[int][]*oai.SchemaRef{}
			for _, errShapeSpec := range opSpec.Errors {
				if errShapeSpec.ShapeName == nil {
					return nil, fmt.Errorf("expected to find error shape name but found %v", errShapeSpec)
				}
				errShapeName := *errShapeSpec.ShapeName
				errShape, found := apiSpec.Shapes[errShapeName]
				if !found {
					return nil, fmt.Errorf("expected to find error shape %s", errShapeName)
				}
				errRespCode := 400
				errShapeSchemaRef := oai.NewSchemaRef("#/components/schemas/"+errShapeName, nil)
				if errShape.Error == nil {
					// Some older XML APIs like S3 do not have an Error field
					// in the shape spec. Instead, the shape spec will have no
					// members, no HTTP response code, nothing but the name of
					// the error.
					//
					// For example:
					//
					// "NoSuchKey":{
					//   "type":"structure",
					//   "members":{
					//   },
					//   "exception":true
					// },
					//
					// In these cases, we just need to use a generic 400 HTTP
					// status code (even though many are actually 404s)
				} else if errShape.Error.HTTPStatusCode != nil {
					errRespCode = *errShape.Error.HTTPStatusCode
				}
				schemasWithCode, exists := codeSchemaMap[errRespCode]
				if !exists {
					schemasWithCode = []*oai.SchemaRef{}
					codeSchemaMap[errRespCode] = schemasWithCode
				}
				codeSchemaMap[errRespCode] = append(codeSchemaMap[errRespCode], errShapeSchemaRef)
			}
			for errRespCode, schemaRefs := range codeSchemaMap {
				if len(schemaRefs) > 1 {
					respSchema := &oai.Schema{OneOf: schemaRefs}
					op.AddResponse(errRespCode, oai.NewResponse().WithJSONSchema(respSchema))
				} else {
					op.AddResponse(errRespCode, oai.NewResponse().WithJSONSchemaRef(schemaRefs[0]))
				}
			}
		}
	}
	return op, nil
}
