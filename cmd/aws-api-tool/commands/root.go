package commands

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	appName      = "aws-api-tool"
	appShortDesc = "aws-api-tool - transform and manipulate AWS API definitions"
	appLongDesc  = `aws-api-tool

A tool to manipulate, transform and normalize AWS API definitions`
)

var (
	version          string
	buildHash        string
	buildDate        string
	debug            bool
	defaultCachePath string
	cachePath        string
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appShortDesc,
	Long:  appLongDesc,
	Args:  processRootCmdArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	hd, err := homedir.Dir()
	if err != nil {
		fmt.Printf("unable to determine $HOME: %s\n", err)
		os.Exit(1)
	}
	defaultCachePath = filepath.Join(hd, "."+appName)
	rootCmd.PersistentFlags().StringVar(
		&cachePath, "cache-path", defaultCachePath, "Path to cache directory root",
	)
	rootCmd.PersistentFlags().BoolVar(
		&debug, "debug", false, "Enable or disable debug mode",
	)
}

func processRootCmdArgs(cmd *cobra.Command, args []string) error {
	if err := processCachePath(); err != nil {
		return err
	}
	return nil
}

func processCachePath() error {
	return os.MkdirAll(cachePath, os.ModePerm)
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute(v string, bh string, bd string) {
	version = v
	buildHash = bh
	buildDate = bd

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
