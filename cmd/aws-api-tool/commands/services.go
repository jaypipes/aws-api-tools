//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/jaypipes/aws-api-tools/pkg/apimodel"
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

func ensureSDKRepo() (string, error) {
	srcPath := filepath.Join(cachePath, "src")
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return "", err
	}
	// clone the aws-sdk-go repository locally so we can query for service
	// information in the models/apis/ directories
	trace("cloning aws-sdk-go to local cache %s ...\n", srcPath)
	clonePath, err := cloneSDKRepo(srcPath)
	if err != nil {
		return "", err
	}
	return clonePath, nil
}

func serviceList(cmd *cobra.Command, args []string) error {
	sdkPath, err := ensureSDKRepo()
	if err != nil {
		return err
	}
	trace("fetching service information from aws-sdk-go ... \n")
	svcs, err := getServices(sdkPath, filteredServices)
	if err != nil {
		return err
	}
	headers := []string{"Alias", "API Version", "Full Name"}
	rows := make([][]string, len(svcs))
	for x, svc := range svcs {
		rows[x] = []string{
			svc.Alias,
			svc.API.Metadata.APIVersion,
			svc.API.Metadata.ServiceFullName,
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.Render()
	return nil
}

type Service struct {
	Alias string
	API   *apimodel.API
}

// getServices returns a slice of Service objects representing the AWS service
// APIs listed in the models/apis/ directory of the aws-sdk-go repository
func getServices(
	clonePath string,
	filtered []string,
) ([]Service, error) {
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
		if !inStrings(fname, filtered) {
			continue
		}
		version, err := getServiceAPIVersion(fp)
		if err != nil {
			return svcs, err
		}
		versionPath := filepath.Join(fp, version)
		api, err := getServiceAPI(versionPath)
		if err != nil {
			return svcs, err
		}
		svcs = append(svcs, Service{Alias: fname, API: api})
	}
	return svcs, nil
}

// getService returns a Service object representing a specified AWS service
func getService(
	clonePath string,
	serviceAlias string,
) (Service, error) {
	svcs, err := getServices(clonePath, []string{serviceAlias})
	if err != nil {
		return Service{}, err
	}
	if len(svcs) == 0 {
		return Service{}, fmt.Errorf("unknown service %s", serviceAlias)
	}
	return svcs[0], nil
}

func getServiceAPIVersion(servicePath string) (string, error) {
	versionDirs, err := ioutil.ReadDir(servicePath)
	if err != nil {
		return "", err
	}
	for _, f := range versionDirs {
		version := f.Name()
		fp := filepath.Join(servicePath, version)
		fi, err := os.Lstat(fp)
		if err != nil {
			return "", err
		}
		if !fi.IsDir() {
			return "", fmt.Errorf("expected to find only directories in service model directory %s but found non-directory %s", servicePath, fp)
		}
		// TODO(jaypipes): handle more than one version? doesn't seem like
		// there is ever more than one.
		return version, nil
	}
	return "", fmt.Errorf("expected to find at least one directory in service model directory %s", servicePath)
}

func getServiceAPI(versionPath string) (*apimodel.API, error) {
	// in each models/apis/$service/$version/ directory will exist files like
	// api-2.json, docs-2.json, etc. We want to grab the API model from the
	// api-2.json file
	modelPath := filepath.Join(versionPath, "api-2.json")
	return apimodel.ParseFrom(modelPath)
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

func inStrings(subject string, collection []string) bool {
	for _, s := range collection {
		if s == subject {
			return true
		}
	}
	return false
}
