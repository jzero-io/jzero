package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	GenModule bool
}

func (c *Config) Wd() string {
	wd, _ := os.Getwd()
	return wd
}

func (c *Config) ProtoDir() string {
	return filepath.Join("desc", "proto")
}

func (c *Config) ApiDir() string {
	return filepath.Join("desc", "api")
}

func (c *Config) SqlDir() string {
	return filepath.Join("desc", "sql")
}
