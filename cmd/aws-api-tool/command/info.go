//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"fmt"

	"github.com/jaypipes/aws-api-tools/pkg/apimodel"
	"github.com/spf13/cobra"
)

// infoCmd lists all object types for an AWS API service
var infoCmd = &cobra.Command{
	Use:   "info <api>",
	Short: "show summary information about an AWS service API",
	Args:  requireAPIArg,
	RunE:  infoAPI,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func infoAPI(cmd *cobra.Command, args []string) error {
	api, err := getAPI(args[0])
	if err != nil {
		return err
	}
	ops := api.GetOperations(nil)
	objects := api.GetObjects(nil)
	countScalars := 0
	countPayloads := 0
	countExceptions := 0
	countLists := 0
	for _, obj := range objects {
		switch obj.Type {
		case apimodel.ObjectTypeScalar:
			countScalars++
		case apimodel.ObjectTypeList:
			countLists++
		case apimodel.ObjectTypePayload:
			countPayloads++
		case apimodel.ObjectTypeException:
			countExceptions++
		}
	}
	fmt.Printf("Full name:        %s\n", api.FullName)
	fmt.Printf("API version:      %s\n", api.Version)
	fmt.Printf("Protocol:         %s\n", api.Protocol)
	fmt.Printf("Total operations: %d\n", len(ops))
	fmt.Printf("Total objects:    %d\n", len(objects))
	fmt.Printf("Total scalars:    %d\n", countScalars)
	fmt.Printf("Total payloads:   %d\n", countPayloads)
	fmt.Printf("Total exceptions: %d\n", countExceptions)
	fmt.Printf("Total lists:      %d\n", countLists)
	return nil
}
