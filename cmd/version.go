/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

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
	Short: "jzero version",
	Long:  `jzero version`,
	RunE:  getVersion,
}

func getVersion(_ *cobra.Command, _ []string) error {
	var versionBuffer bytes.Buffer

	versionBuffer.WriteString(fmt.Sprintf("jzero version %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH))
	versionBuffer.WriteString(fmt.Sprintf("Go version %s\n", runtime.Version()))
	if Commit != "" {
		versionBuffer.WriteString(fmt.Sprintf("Git commit %s\n", Commit))
	}
	if Date != "" {
		versionBuffer.WriteString(fmt.Sprintf("Build date %s\n", cast.ToTimeInDefaultLocation(Date, time.Local)))
	}

	fmt.Print(versionBuffer.String())
	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
