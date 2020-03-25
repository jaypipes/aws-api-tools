package main

import (
	"github.com/jaypipes/aws-api-tools/cmd/aws-api-tool/commands"
)

var (
	// version of application at compile time (-X 'main.version=$(VERSION)').
	version = "(Unknown Version)"
	// buildHash GIT hash of application at compile time (-X 'main.buildHash=$(GITCOMMIT)').
	buildHash = "No Git-hash Provided."
	// buildDate of application at compile time (-X 'main.buildDate=$(BUILDDATE)').
	buildDate = "No Build Date Provided."
)

func main() {
	commands.Execute(version, buildHash, buildDate)
}
