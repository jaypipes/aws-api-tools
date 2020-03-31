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
	Method     string `json:"method"`
	RequestURI string `json:"requestUri"`
}

type opSpec struct {
	HTTP   httpSpec       `json:"http"`
	Input  shapeRefSpec   `json:"input"`
	Output shapeRefSpec   `json:"output"`
	Errors []shapeRefSpec `json:"errors"`
}

type shapeSpec struct {
	Type       string                  `json:"type"`
	Exception  bool                    `json:"exception"`
	Required   []string                `json:"required"`
	Members    map[string]shapeRefSpec `json:"members"`
	ListMember *shapeRefSpec           `json:"member",omitempty` // for list types
}

type apiSpec struct {
	Metadata   metadataSpec         `json:"metadata"`
	Operations map[string]opSpec    `json:"operations"`
	Shapes     map[string]shapeSpec `json:"shapes"`
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
	if len(api.objectMap) > 0 {
		return nil
	}
	spec := api.apiSpec
	api.objectMap = make(map[string]*Object, len(spec.Shapes))
	api.opMap = make(map[string]Operation, len(spec.Operations))
	api.payloadMap = map[string]*Object{}
	api.scalarMap = map[string]*Object{}
	api.exceptionMap = map[string]*Object{}
	api.listMap = map[string]*Object{}
	api.resourceMap = map[string]*Resource{}

	// Populate the base object maps
	for shapeName, shapeSpec := range spec.Shapes {
		obj := Object{
			Name:                shapeName,
			Type:                ObjectTypeObject,
			DataType:            shapeSpec.Type,
			Members:             make(map[string]*Object, len(shapeSpec.Members)),
			RequiredMemberNames: shapeSpec.Required,
		}
		// Determine simple types like scalars, lists and exceptions
		if shapeSpec.Type != "structure" && shapeSpec.Type != "list" {
			obj.Type = ObjectTypeScalar
			api.scalarMap[shapeName] = &obj
		} else if shapeSpec.Type == "structure" && shapeSpec.Exception {
			obj.Type = ObjectTypeException
			api.exceptionMap[shapeName] = &obj
		} else if shapeSpec.Type == "list" {
			obj.Type = ObjectTypeList
			api.listMap[shapeName] = &obj
		}
		api.objectMap[shapeName] = &obj

	}

	// Set each object's member references
	for shapeName, shapeSpec := range spec.Shapes {
		obj, found := api.objectMap[shapeName]
		if !found {
			return fmt.Errorf("expected to find object %s in objectMap", shapeName)
		}
		// List types are special...
		if shapeSpec.ListMember != nil {
			listMemberShapeName := *shapeSpec.ListMember.ShapeName
			listMemberObj, found := api.objectMap[listMemberShapeName]
			if !found {
				return fmt.Errorf("expected to find member object %s in objectMap", listMemberShapeName)
			}
			obj.Members[listMemberShapeName] = listMemberObj
			continue
		}
		if len(shapeSpec.Members) == 0 {
			continue
		}
		x := 0
		for memberName, memberShapeRef := range shapeSpec.Members {
			if memberShapeRef.ShapeName == nil {
				continue
			}
			memberShapeName := *memberShapeRef.ShapeName
			memberObj, found := api.objectMap[memberShapeName]
			if !found {
				return fmt.Errorf("expected to find member object %s in objectMap", memberShapeName)
			}
			obj.Members[memberName] = memberObj
			x++
		}
	}
	for opName, opSpec := range spec.Operations {
		op := Operation{
			Name:       opName,
			Method:     opSpec.HTTP.Method,
			RequestURI: opSpec.HTTP.RequestURI,
		}
		if opSpec.Input.ShapeName != nil {
			inShapeName := *opSpec.Input.ShapeName
			inObj, ok := api.objectMap[inShapeName]
			if !ok {
				return fmt.Errorf("expected to find object %s", inShapeName)
			}
			inObj.Type = ObjectTypePayload
			api.payloadMap[inShapeName] = inObj
			op.Input = inObj
		}
		if opSpec.Output.ShapeName != nil {
			outShapeName := *opSpec.Output.ShapeName
			outObj, ok := api.objectMap[outShapeName]
			if !ok {
				return fmt.Errorf("expected to find object %s", outShapeName)
			}
			outObj.Type = ObjectTypePayload
			api.payloadMap[outShapeName] = outObj
			op.Output = outObj
		}
		if len(opSpec.Errors) > 0 {
			errs := make([]*Object, len(opSpec.Errors))
			for x, errShapeRef := range opSpec.Errors {
				errShapeName := *errShapeRef.ShapeName
				errObj, ok := api.objectMap[errShapeName]
				if !ok {
					return fmt.Errorf("expected to find object %s", errShapeName)
				}
				errs[x] = errObj
			}
			op.Errors = errs
		}
		api.opMap[opName] = op
	}

	resources, err := getResources(api)
	if err != nil {
		return err
	}
	api.resourceMap = resources
	return nil
}
