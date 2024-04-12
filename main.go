package main

import (
	"embed"
	"github.com/jaronnie/jzero/cmd"
	"github.com/jaronnie/jzero/embedx"
)

//go:embed .template
var template embed.FS

//go:embed .protosets/*.pb
var protosets embed.FS

//go:embed config.toml
var config embed.FS

//go:embed all:web
var web embed.FS

func main() {
	embedx.Web = web
	embedx.Protosets = protosets
	embedx.Config = config
	embedx.Template = template
	cmd.Execute()
}
