//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var (
	cliOutputFormat string
)

// schemaCmd shows a schema document for an AWS API service
var schemaCmd = &cobra.Command{
	Use:   "schema <api>",
	Short: "shows OpenAPI schema information for an AWS service API",
	Args:  requireAPIArg,
	RunE:  apiSchema,
}

func init() {
	schemaCmd.PersistentFlags().StringVarP(
		&cliOutputFormat, "format", "f", "yaml", "Output format for schema ('yaml' or 'json').",
	)
	rootCmd.AddCommand(schemaCmd)
}

func apiSchema(cmd *cobra.Command, args []string) error {
	api, err := getAPI(args[0])
	if err != nil {
		return err
	}
	json, err := api.Schema().MarshalJSON()
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
