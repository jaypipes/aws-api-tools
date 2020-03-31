//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var (
	cliOutputFormat string
)
var requireAPIAndResourceArg = func(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("requires an <api> and <resource> argument")
	}
	return nil
}

// schemaCmd shows a schema document for an AWS API service and resource
var schemaCmd = &cobra.Command{
	Use:   "schema <api> <resource>",
	Short: "shows schema information for an AWS service API and resource",
	Args:  requireAPIAndResourceArg,
	RunE:  resourceSchema,
}

func init() {
	schemaCmd.PersistentFlags().StringVarP(
		&cliOutputFormat, "format", "f", "yaml", "Output format for schema ('yaml' or 'json').",
	)
	rootCmd.AddCommand(schemaCmd)
}

func resourceSchema(cmd *cobra.Command, args []string) error {
	api, err := getAPI(args[0])
	if err != nil {
		return err
	}
	resource, err := api.GetResource(args[1])
	if err != nil {
		return err
	}
	json, err := resource.OpenAPI3Schema().MarshalJSON()
	if err != nil {
		return err
	}
	if cliOutputFormat == "yaml" {
		yamlStr, err := yaml.JSONToYAML(json)
		if err != nil {
			return err
		}
		fmt.Printf(string(yamlStr))
	} else {
		fmt.Println(string(json))
	}
	return nil
}
