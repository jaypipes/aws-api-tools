//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	sdkRepoURL = "https://github.com/aws/aws-sdk-go"
)

var (
	filteredServices = []string{}
	cliServices      = ""
)

// servicesCmd provides sub-commands for querying/listing AWS services
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "query names and other information about AWS service APIs",
	Args:  processServiceCmdArgs,
}

// servicesListCmd lists AWS service APIs
var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists AWS service API information",
	Args:  processServiceCmdArgs,
	RunE:  serviceList,
}

func init() {
	serviceCmd.PersistentFlags().StringVar(
		&cliServices, "services", "", "Comma-delimited list of AWS service aliases to operate on (e.g. --services s3,iam) Default is to operate on all services.",
	)
	serviceCmd.AddCommand(serviceListCmd)
	rootCmd.AddCommand(serviceCmd)
}

func processServiceCmdArgs(cmd *cobra.Command, args []string) error {
	processFilteredServices()
	return nil
}

func processFilteredServices() {
	if cliServices != "" {
		filteredServices = strings.Split(cliServices, ",")
	}
}

func serviceList(cmd *cobra.Command, args []string) error {
	srcPath := filepath.Join(cachePath, "src")
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}
	// clone the aws-sdk-go repository locally so we can query for service
	// information in the models/apis/ directories
	clonePath, err := cloneSDKRepo(srcPath)
	if err != nil {
		return err
	}
	svcs, err := getServices(clonePath, filteredServices)
	if err != nil {
		return err
	}
	headers := []string{"Name"}
	rows := make([][]string, len(svcs))
	for x, svc := range svcs {
		rows[x] = []string{svc.Alias}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
	return nil
}

type Service struct {
	Alias         string
	Versions      []string
	LatestVersion string
}

func getServices(clonePath string, filteredServices []string) ([]Service, error) {
	svcs := []Service{}

	destPath := filepath.Join(clonePath, "models", "apis")
	apiDirs, err := ioutil.ReadDir(destPath)
	if err != nil {
		return svcs, err
	}
	for _, f := range apiDirs {
		fname := f.Name()
		fp := filepath.Join(destPath, fname)
		fi, err := os.Lstat(fp)
		if err != nil {
			return svcs, err
		}
		if !fi.IsDir() {
			continue
		}
		// Filter just the services we're interested in
		if cliServices != "" {
			if !inFilteredServices(fname) {
				continue
			}
		}
		svcs = append(svcs, Service{Alias: fname})
	}
	return svcs, nil
}

// cloneSDKRepo git clone's the aws-sdk-go source repo into the cache and
// returns the filepath to the clone'd repo
func cloneSDKRepo(srcPath string) (string, error) {
	clonePath := filepath.Join(srcPath, "aws-sdk-go")
	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "--depth", "1", sdkRepoURL, clonePath)
		return clonePath, cmd.Run()
	}
	return clonePath, nil
}

func inFilteredServices(service string) bool {
	for _, s := range filteredServices {
		if s == service {
			return true
		}
	}
	return false
}
