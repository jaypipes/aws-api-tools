//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jaypipes/aws-api-tools/pkg/apimodel"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cliService               string
	cliHTTPMethodFilter      string
	cliOperationPrefixFilter string
	serviceRef               *Service
)

// apiCmd provides sub-commands for exploring the AWS API models
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "query API model information for an AWS service API",
}

// apiInfoCmd shows summary information for one or more AWS API service models
var apiInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "shows summary information for an AWS service API",
	RunE:  apiInfo,
}

// apiOperationsCmd lists all operations for an AWS API service
var apiOperationsCmd = &cobra.Command{
	Use:   "operations",
	Short: "lists Operations for an AWS API service",
	RunE:  apiOperations,
}

// apiPrimariesCmd lists all primary object types for an AWS API service
var apiPrimariesCmd = &cobra.Command{
	Use:   "primaries",
	Short: "lists Primary object types for an AWS API service",
	RunE:  apiPrimaries,
}

// apiObjectsCmd lists all object types for an AWS API service
var apiObjectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "lists Object types for an AWS API service",
	RunE:  apiObjects,
}

// apiScalarsCmd lists all scalar types for an AWS API service
var apiScalarsCmd = &cobra.Command{
	Use:   "scalars",
	Short: "lists Scalar types for an AWS API service",
	RunE:  apiScalars,
}

// apiPayloadsCmd lists all payload types for an AWS API service
var apiPayloadsCmd = &cobra.Command{
	Use:   "payloads",
	Short: "lists Payload types for an AWS API service",
	RunE:  apiPayloads,
}

// apiExceptionsCmd lists all exception types for an AWS API service
var apiExceptionsCmd = &cobra.Command{
	Use:   "exceptions",
	Short: "lists Exception types for an AWS API service",
	RunE:  apiExceptions,
}

// apiListsCmd lists all list types for an AWS API service
var apiListsCmd = &cobra.Command{
	Use:   "lists",
	Short: "lists List types for an AWS API service",
	RunE:  apiLists,
}

func init() {
	apiCmd.PersistentFlags().StringVarP(
		&cliService, "service", "s", "", "Alias of the AWS service to work with.",
	)
	apiOperationsCmd.PersistentFlags().StringVar(
		&cliOperationPrefixFilter, "prefix", "", "Comma-delimited list of string prefixes to filter operations by.",
	)
	apiOperationsCmd.PersistentFlags().StringVarP(
		&cliHTTPMethodFilter, "method", "m", "", "Comma-delimited list of HTTP methods to filter operations by.",
	)
	apiCmd.MarkFlagRequired("service")
	apiCmd.AddCommand(apiInfoCmd)
	apiCmd.AddCommand(apiOperationsCmd)
	apiCmd.AddCommand(apiPrimariesCmd)
	apiCmd.AddCommand(apiObjectsCmd)
	apiCmd.AddCommand(apiScalarsCmd)
	apiCmd.AddCommand(apiPayloadsCmd)
	apiCmd.AddCommand(apiExceptionsCmd)
	apiCmd.AddCommand(apiListsCmd)
	rootCmd.AddCommand(apiCmd)
}

func ensureService() error {
	sdkPath, err := ensureSDKRepo()
	if err != nil {
		return err
	}
	trace("fetching service information from aws-sdk-go ... \n")
	svc, err := getService(sdkPath, cliService)
	if err != nil {
		return err
	}
	serviceRef = &svc
	return nil
}

func apiInfo(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIInfo(serviceRef)
	return nil
}

func printAPIInfo(svc *Service) {
	ops := svc.API.GetOperations(nil)
	objects := svc.API.GetObjects()
	scalars := svc.API.GetScalars()
	payloads := svc.API.GetPayloads()
	exceptions := svc.API.GetExceptions()
	lists := svc.API.GetLists()
	fmt.Printf("Service '%s'\n", svc.Alias)
	fmt.Printf("  API Version:         %s\n", svc.API.Metadata.APIVersion)
	fmt.Printf("  Total operations:    %d\n", len(ops))
	fmt.Printf("  Total scalars:       %d\n", len(scalars))
	fmt.Printf("  Total objects:       %d\n", len(objects))
	fmt.Printf("  Total payloads:      %d\n", len(payloads))
	fmt.Printf("  Total exceptions:    %d\n", len(exceptions))
	fmt.Printf("  Total lists:         %d\n", len(lists))
}

func apiOperations(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIOperations(serviceRef)
	return nil
}

func printAPIOperations(svc *Service) {
	filter := &apimodel.OperationFilter{}
	if cliHTTPMethodFilter != "" {
		filter.Methods = strings.Split(strings.ToUpper(cliHTTPMethodFilter), ",")

	}
	if cliOperationPrefixFilter != "" {
		filter.Prefixes = strings.Split(cliOperationPrefixFilter, ",")
	}
	operations := svc.API.GetOperations(filter)
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
}

func apiPrimaries(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIPrimaries(serviceRef)
	return nil
}

func printAPIPrimaries(svc *Service) {
	primaries := svc.API.GetPrimaries()
	headers := []string{"Name"}
	rows := make([][]string, len(primaries))
	for x, primary := range primaries {
		rows[x] = []string{primary.SingularName}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

func apiObjects(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIObjects(serviceRef)
	return nil
}

func printAPIObjects(svc *Service) {
	objects := svc.API.GetObjects()
	headers := []string{"Name"}
	rows := make([][]string, len(objects))
	for x, object := range objects {
		rows[x] = []string{object.Name}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

func apiScalars(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIScalars(serviceRef)
	return nil
}

func printAPIScalars(svc *Service) {
	scalars := svc.API.GetScalars()
	headers := []string{"Name", "Type"}
	rows := make([][]string, len(scalars))
	for x, scalar := range scalars {
		rows[x] = []string{scalar.Name, scalar.Type}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

func apiPayloads(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIPayloads(serviceRef)
	return nil
}

func printAPIPayloads(svc *Service) {
	payloads := svc.API.GetPayloads()
	headers := []string{"Name"}
	rows := make([][]string, len(payloads))
	for x, payload := range payloads {
		rows[x] = []string{payload.Name}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

func apiExceptions(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIExceptions(serviceRef)
	return nil
}

func printAPIExceptions(svc *Service) {
	exceptions := svc.API.GetExceptions()
	headers := []string{"Name"}
	rows := make([][]string, len(exceptions))
	for x, exception := range exceptions {
		rows[x] = []string{exception.Name}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}

func apiLists(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPILists(serviceRef)
	return nil
}

func printAPILists(svc *Service) {
	lists := svc.API.GetLists()
	headers := []string{"Name"}
	rows := make([][]string, len(lists))
	for x, list := range lists {
		rows[x] = []string{list.Name}
	}
	noResults(rows)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}
