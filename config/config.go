package config

import "embed"

//go:embed config.toml
var Config embed.FS
