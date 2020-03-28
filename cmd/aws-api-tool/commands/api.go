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
	Short: "shows summary information for one or more AWS API service models",
	RunE:  apiInfo,
}

// apiOperationsCmd lists all operations for an AWS API service
var apiOperationsCmd = &cobra.Command{
	Use:   "operations",
	Short: "lists operations for an AWS API service",
	RunE:  apiOperations,
}

// apiObjectsCmd lists all scalar types for an AWS API service
var apiObjectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "lists object types for an AWS API service",
	RunE:  apiObjects,
}

// apiScalarsCmd lists all scalar types for an AWS API service
var apiScalarsCmd = &cobra.Command{
	Use:   "scalars",
	Short: "lists scalar types for an AWS API service",
	RunE:  apiScalars,
}

// apiPayloadsCmd lists all scalar types for an AWS API service
var apiPayloadsCmd = &cobra.Command{
	Use:   "payloads",
	Short: "lists payload types for an AWS API service",
	RunE:  apiPayloads,
}

// apiExceptionsCmd lists all scalar types for an AWS API service
var apiExceptionsCmd = &cobra.Command{
	Use:   "exceptions",
	Short: "lists exception types for an AWS API service",
	RunE:  apiExceptions,
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
	apiCmd.AddCommand(apiObjectsCmd)
	apiCmd.AddCommand(apiScalarsCmd)
	apiCmd.AddCommand(apiPayloadsCmd)
	apiCmd.AddCommand(apiExceptionsCmd)
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
	objects := svc.API.GetObjects()
	scalars := svc.API.GetScalars()
	payloads := svc.API.GetPayloads()
	exceptions := svc.API.GetExceptions()
	fmt.Printf("Service '%s'\n", svc.Alias)
	fmt.Printf("  API Version:         %s\n", svc.API.Metadata.APIVersion)
	fmt.Printf("  Total operations:    %d\n", len(svc.API.Operations))
	fmt.Printf("  Total shapes:        %d\n", len(svc.API.Shapes))
	fmt.Printf("  Total scalars:       %d\n", len(scalars))
	fmt.Printf("  Total objects:       %d\n", len(objects))
	fmt.Printf("  Total payloads:      %d\n", len(payloads))
	fmt.Printf("  Total exceptions:    %d\n", len(exceptions))
}

func apiOperations(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIOperations(serviceRef)
	return nil
}

func printAPIOperations(svc *Service) {
	filterMethods := strings.Split(strings.ToUpper(cliHTTPMethodFilter), ",")
	filterPrefixes := strings.Split(cliOperationPrefixFilter, ",")
	operations := svc.API.GetOperations(filterMethods, filterPrefixes)
	headers := []string{"Name", "HTTP Method"}
	rows := make([][]string, 0)
	for operationName, operation := range operations {
		rows = append(rows, []string{operationName, operation.HTTP.Method})
	}
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
	x := 0
	for objectName, _ := range objects {
		rows[x] = []string{objectName}
		x++
	}
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
	x := 0
	for scalarName, scalarType := range scalars {
		rows[x] = []string{scalarName, scalarType}
		x++
	}
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
	x := 0
	for payloadName, _ := range payloads {
		rows[x] = []string{payloadName}
		x++
	}
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
	x := 0
	for exceptionName, _ := range exceptions {
		rows[x] = []string{exceptionName}
		x++
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
}
