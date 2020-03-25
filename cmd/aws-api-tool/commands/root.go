package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	appName      = "aws-api-tool"
	appShortDesc = "aws-api-tool - transform and manipulate AWS API definitions"
	appLongDesc  = `aws-api-tool

A tool to manipulate, transform and normalize AWS API definitions`
)

var (
	version   string
	buildHash string
	buildDate string
	debug     bool
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appShortDesc,
	Long:  appLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
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

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&debug, "debug", false, "Enable or disable debug mode",
	)
}
