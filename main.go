package main

import (
	"embed"
	"os"
	"strings"

	"github.com/jzero-io/jzero/cmd"
	"github.com/jzero-io/jzero/embeded"
)

// embeded
var (
	//go:embed all:.template
	template embed.FS
)

// ldflags
var (
	version = "0.38.0-alpha"
	commit  string
	date    string
)

func main() {
	embeded.Template = template
	cmd.Version = version
	cmd.Date = date
	cmd.Commit = commit

	args := os.Args[1:]

	mcpArgs := []string{os.Args[0]}
	for _, a := range args {
		split := strings.Split(a, "__")
		mcpArgs = append(mcpArgs, split...)
	}

	os.Args = mcpArgs

	cmd.Execute()
}
