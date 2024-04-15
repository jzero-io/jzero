package main

import (
	"embed"
	"github.com/jaronnie/jzero/cmd"
	"github.com/jaronnie/jzero/embeded"
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
	embeded.Web = web
	embeded.Protosets = protosets
	embeded.Config = config
	embeded.Template = template
	cmd.Execute()
}
