//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package model

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrInvalidVersionDirectory = errors.New(
		"expected to find only directories in api model directory but found non-directory",
	)
	ErrNoValidVersionDirectory = errors.New(
		"no valid version directories found",
	)
)

// SDKHelper is a helper struct that helps work with the aws-sdk-go models and
// API model loader
type SDKHelper struct {
	basePath string
}

// NewSDKHelper returns a new SDKHelper object
func NewSDKHelper(basePath string) *SDKHelper {
	return &SDKHelper{basePath}
}

// APIVersion returns the API version (e.g. "2012-10-03") for a service API
func (h *SDKHelper) APIVersion(serviceAlias string) (string, error) {
	apiPath := filepath.Join(h.basePath, "models", "apis", serviceAlias)
	versionDirs, err := ioutil.ReadDir(apiPath)
	if err != nil {
		return "", err
	}
	for _, f := range versionDirs {
		version := f.Name()
		fp := filepath.Join(apiPath, version)
		fi, err := os.Lstat(fp)
		if err != nil {
			return "", err
		}
		if !fi.IsDir() {
			return "", ErrInvalidVersionDirectory
		}
		// TODO(jaypipes): handle more than one version? doesn't seem like
		// there is ever more than one.
		return version, nil
	}
	return "", ErrNoValidVersionDirectory
}
