package protosets

import "embed"

//go:embed *.pb
var Protosets embed.FS
