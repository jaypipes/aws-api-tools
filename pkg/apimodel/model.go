//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package apimodel

type Metadata struct {
	APIVersion      string `json:"apiVersion"`
	ServiceFullName string `json:"serviceFullName"`
}

type Shape struct {
	Type      string `json:"type"`
	Exception bool   `json:"exception"`
}

type API struct {
	Metadata Metadata         `json:"metadata"`
	Shapes   map[string]Shape `json:"shapes"`
}

// Scalars returns the Shapes in the API that are of a non-compound type by
// returning a map of the shape name and its underlying simple type
func (a API) Scalars() map[string]string {
	res := map[string]string{}
	for shapeName, shape := range a.Shapes {
		if shape.Type != "structure" && shape.Type != "list" {
			res[shapeName] = shape.Type
		}
	}
	return res
}
