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
	cliResource     string
	cliOutputFormat string
)

// resourceCmd provides sub-commands for exploring a specific resource in an
// AWS service API
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "query and transform information about a specific resource in an AWS service API",
}

// resourceSchemaCmd shows an OpenAPI3 Schema for a specific resource in an AWS
// service APi
var openapi3SchemaCmd = &cobra.Command{
	Use:     "openapi3",
	Aliases: []string{"oai3", "openapi"},
	Short:   "display OpenAPI3 Schema for a specific resource in an AWS service API",
	RunE:    openapi3Schema,
}

func init() {
	schemaCmd.PersistentFlags().StringVarP(
		&cliOutputFormat, "format", "f", "yaml", "Output format for schema ('yaml' or 'json').",
	)
	schemaCmd.PersistentFlags().StringVarP(
		&cliService, "service", "s", "", "Alias of the AWS service to work with.",
	)
	schemaCmd.MarkFlagRequired("service")
	schemaCmd.PersistentFlags().StringVarP(
		&cliResource, "resource", "r", "", "Name of the API resource to work with.",
	)
	schemaCmd.MarkFlagRequired("resource")
	schemaCmd.AddCommand(openapi3SchemaCmd)
	rootCmd.AddCommand(schemaCmd)
}

func openapi3Schema(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	return printOpenAPI3Schema(serviceRef, cliResource)
}

func printOpenAPI3Schema(svc *Service, resourceName string) error {
	resource, err := svc.API.GetResource(resourceName)
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
