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

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	cliService string
	serviceRef *Service
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

// apiScalarsCmd lists all scalar types for an AWS API service
var apiScalarsCmd = &cobra.Command{
	Use:   "scalars",
	Short: "lists scalar types for an AWS API service",
	RunE:  apiScalars,
}

func init() {
	apiCmd.PersistentFlags().StringVar(
		&cliService, "service", "s", "Alias of the AWS service to work with.",
	)
	apiCmd.MarkFlagRequired("service")
	apiCmd.AddCommand(apiInfoCmd)
	apiCmd.AddCommand(apiScalarsCmd)
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
	scalars := svc.API.Scalars()
	fmt.Printf("Service '%s'\n", svc.Alias)
	fmt.Printf("  API Version:         %s\n", svc.API.Metadata.APIVersion)
	fmt.Printf("  Total shapes:        %d\n", len(svc.API.Shapes))
	fmt.Printf("  Total scalars:       %d\n", len(scalars))
}

func apiScalars(cmd *cobra.Command, args []string) error {
	if err := ensureService(); err != nil {
		return err
	}
	printAPIScalars(serviceRef)
	return nil
}

func printAPIScalars(svc *Service) {
	scalars := svc.API.Scalars()
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
