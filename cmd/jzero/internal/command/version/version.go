/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package version

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: `Print jzero version`,
	Run: func(cmd *cobra.Command, args []string) {
		GetVersion()
	},
}

func GetVersion() {
	var versionBuffer bytes.Buffer

	if Version != "" {
		versionBuffer.WriteString(fmt.Sprintf("jzero version %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH))
	} else {
		versionBuffer.WriteString(fmt.Sprintf("jzero version %s %s/%s\n", "unknown", runtime.GOOS, runtime.GOARCH))
	}

	versionBuffer.WriteString(fmt.Sprintf("Go version %s\n", runtime.Version()))
	if Commit != "" {
		versionBuffer.WriteString(fmt.Sprintf("Git commit %s\n", Commit))
	} else {
		versionBuffer.WriteString(fmt.Sprintf("Git commit %s\n", "unknown"))
	}

	if Date != "" {
		versionBuffer.WriteString(fmt.Sprintf("Build date %s\n", cast.ToTimeInDefaultLocation(Date, time.Local)))
	} else {
		versionBuffer.WriteString(fmt.Sprintf("Build date %s\n", "unknown"))
	}

	fmt.Print(versionBuffer.String())
}

func GetCommand() *cobra.Command {
	return versionCmd
}
