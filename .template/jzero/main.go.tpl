package main

import (
	"{{ .Module }}/cmd"
)

// ldflags
var (
	version string
	commit  string
	date    string
)

func main() {
    cmd.Version = version
    cmd.Date = date
    cmd.Commit = commit

	cmd.Execute()
}