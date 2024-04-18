/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
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

func getVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf(`Jzero version: %s
Go version: %s
Commit: %s
Date: %s
`, Version, runtime.Version(), Commit, cast.ToTimeInDefaultLocation(Date, time.Local))
	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
