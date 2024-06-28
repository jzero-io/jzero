package main

import (
	"embed"

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
	version = "0.19.1"
	commit  string
	date    string
)

func main() {
	embeded.Template = template
	cmd.Version = version
	cmd.Date = date
	cmd.Commit = commit

	cmd.Execute()
}
