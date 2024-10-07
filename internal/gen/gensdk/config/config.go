package config

import "github.com/jzero-io/jzero/config"

type Config struct {
	GenModule bool
	config.GenConfig
	config.GenSdkConfig
}
