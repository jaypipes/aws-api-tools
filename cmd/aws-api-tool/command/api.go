//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jaypipes/aws-api-tools/pkg/apimodel"
	"github.com/jaypipes/aws-api-tools/pkg/model"
)

const (
	sdkRepoURL = "https://github.com/aws/aws-sdk-go"
)

func ensureSDKRepo() (string, error) {
	srcPath := filepath.Join(cachePath, "src")
	if err := os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return "", err
	}
	// clone the aws-sdk-go repository locally so we can query for API
	// information in the models/apis/ directories
	trace("cloning aws-sdk-go to local cache %s ...\n", srcPath)
	clonePath, err := cloneSDKRepo(srcPath)
	if err != nil {
		return "", err
	}
	return clonePath, nil
}

type APIFilter struct {
	anyMatch         []string
	anyProtocolMatch []string
}

// getAPIs returns a slice of pointer to apimodel.API objects representing the
// AWS service APIs listed in the models/apis/ directory of the aws-sdk-go
// repository
func getAPIs(
	filter *APIFilter,
) ([]*apimodel.API, error) {
	sdkPath, err := ensureSDKRepo()
	if err != nil {
		return nil, err
	}
	sdkHelper := model.NewSDKHelper(sdkPath)
	apis := []*apimodel.API{}

	destPath := filepath.Join(sdkPath, "models", "apis")
	apiDirs, err := ioutil.ReadDir(destPath)
	if err != nil {
		return apis, err
	}
	for _, f := range apiDirs {
		fname := f.Name()
		fp := filepath.Join(destPath, fname)
		fi, err := os.Lstat(fp)
		if err != nil {
			return apis, err
		}
		if !fi.IsDir() {
			continue
		}
		if filter != nil && len(filter.anyMatch) > 0 {
			if !inStrings(fname, filter.anyMatch) {
				continue
			}
		}
		version, err := sdkHelper.APIVersion(fname)
		if err != nil {
			return apis, err
		}
		versionPath := filepath.Join(fp, version)
		api, err := getAPIFromVersionPath(fname, versionPath)
		if err != nil {
			return apis, err
		}
		if filter != nil && len(filter.anyProtocolMatch) > 0 {
			if !inStrings(api.Protocol, filter.anyProtocolMatch) {
				continue
			}
		}
		apis = append(apis, api)
	}
	return apis, nil
}

// getAPI returns a pointer to an apimodel.API object representing a
// specified AWS service
func getAPI(
	alias string,
) (*apimodel.API, error) {
	apis, err := getAPIs(&APIFilter{anyMatch: []string{alias}})
	if err != nil {
		return nil, err
	}
	if len(apis) == 0 {
		return nil, fmt.Errorf("unknown API %s", alias)
	}
	return apis[0], nil
}

func getAPIFromVersionPath(
	alias string,
	versionPath string,
) (*apimodel.API, error) {
	// in each models/apis/$service/$version/ directory will exist files like
	// api-2.json, docs-2.json, etc. We want to grab the API model from the
	// api-2.json file
	modelPath := filepath.Join(versionPath, "api-2.json")
	docPath := filepath.Join(versionPath, "docs-2.json")
	return apimodel.New(alias, modelPath, docPath)
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
	if len(collection) == 0 {
		return true
	}
	for _, s := range collection {
		if s == subject {
			return true
		}
	}
	return false
}
