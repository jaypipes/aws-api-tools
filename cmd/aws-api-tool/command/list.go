//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/jaypipes/aws-api-tools/pkg/apimodel"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cliListAPIsFilter                 string
	cliListOperationsHTTPMethodFilter string
	cliListOperationsPrefixFilter     string
	cliListObjectsTypeFilter          string
	cliListObjectsPrefixFilter        string
)

// listAPIsCmd lists AWS service APIs
var listAPIsCmd = &cobra.Command{
	Use:     "list-apis",
	Aliases: []string{"apis"},
	Short:   "lists AWS service APIs",
	RunE:    listAPIs,
}

var requireAPIArg = func(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires an <api> argument")
	}
	return nil
}

// listOperationsCmd lists all operations for an AWS API service
var listOperationsCmd = &cobra.Command{
	Use:     "list-operations <api>",
	Aliases: []string{"ops"},
	Args:    requireAPIArg,
	Short:   "lists Operations for an AWS service API",
	RunE:    listOperations,
}

// listObjectsCmd lists all object types for an AWS API service
var listObjectsCmd = &cobra.Command{
	Use:     "list-objects <api>",
	Aliases: []string{"objects"},
	Short:   "lists Object types for an AWS service API",
	Args:    requireAPIArg,
	RunE:    listObjects,
}

func init() {
	listAPIsCmd.PersistentFlags().StringVarP(
		&cliListAPIsFilter, "filter", "f", "", "Comma-delimited list of strings to filter APIs on.",
	)
	listOperationsCmd.PersistentFlags().StringVarP(
		&cliListOperationsPrefixFilter, "prefix", "p", "", "Comma-delimited list of string prefixes to filter operations by.",
	)
	listOperationsCmd.PersistentFlags().StringVarP(
		&cliListOperationsHTTPMethodFilter, "method", "m", "", "Comma-delimited list of HTTP methods to filter operations by.",
	)
	listObjectsCmd.PersistentFlags().StringVarP(
		&cliListObjectsPrefixFilter, "prefix", "p", "", "Comma-delimited list of string prefixes to filter objects by.",
	)
	listObjectsCmd.PersistentFlags().StringVarP(
		&cliListObjectsTypeFilter, "type", "t", "", "Comma-delimited list of object types to filter objects by.",
	)
	rootCmd.AddCommand(listAPIsCmd)
	rootCmd.AddCommand(listOperationsCmd)
	rootCmd.AddCommand(listObjectsCmd)
}

func listAPIs(cmd *cobra.Command, args []string) error {
	var filter *APIFilter
	if cliListAPIsFilter != "" {
		filter = &APIFilter{
			anyMatch: strings.Split(cliListAPIsFilter, ","),
		}
	}
	apis, err := getAPIs(filter)
	if err != nil {
		return err
	}
	headers := []string{"Alias", "API Version", "Full Name"}
	rows := make([][]string, len(apis))
	for x, api := range apis {
		rows[x] = []string{
			api.Alias,
			api.Version,
			api.FullName,
		}
	}
	noResults(rows)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func listOperations(cmd *cobra.Command, args []string) error {
	api, err := getAPI(args[0])
	if err != nil {
		return err
	}
	filter := &apimodel.OperationFilter{}
	if cliListOperationsHTTPMethodFilter != "" {
		filter.Methods = strings.Split(
			strings.ToUpper(cliListOperationsHTTPMethodFilter), ",",
		)
	}
	if cliListOperationsPrefixFilter != "" {
		filter.Prefixes = strings.Split(cliListOperationsPrefixFilter, ",")
	}
	operations := api.GetOperations(filter)
	headers := []string{"Name", "HTTP Method"}
	rows := make([][]string, len(operations))
	for x, operation := range operations {
		rows[x] = []string{operation.Name, operation.Method}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func listObjects(cmd *cobra.Command, args []string) error {
	api, err := getAPI(args[0])
	if err != nil {
		return err
	}
	filter := &apimodel.ObjectFilter{}
	if cliListObjectsTypeFilter != "" {
		filter.Types = strings.Split(cliListObjectsTypeFilter, ",")
	}
	if cliListObjectsPrefixFilter != "" {
		filter.Prefixes = strings.Split(cliListObjectsPrefixFilter, ",")
	}
	objects := api.GetObjects(filter)
	headers := []string{"Name", "Object Type", "Data Type"}
	rows := make([][]string, len(objects))
	for x, object := range objects {
		rows[x] = []string{
			object.Name,
			object.Type,
			object.DataType,
		}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
	return nil
}
