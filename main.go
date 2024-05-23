package main

import (
	"embed"

	"github.com/jzero-io/jzero/cmd"
	"github.com/jzero-io/jzero/embeded"
)

// embeded
var (
	//go:embed .template
	template embed.FS

	//go:embed .protosets/*.pb
	protosets embed.FS

	//go:embed config.toml
	config embed.FS

	//go:embed all:web
	web embed.FS
)

// ldflags
var (
	version = "0.17.2"
	commit  string
	date    string
)

func main() {
	{
		embeded.Web = web
		embeded.Protosets = protosets
		embeded.Config = config
		embeded.Template = template
	}

	{
		cmd.Version = version
		cmd.Date = date
		cmd.Commit = commit
	}

	cmd.Execute()
}
