package cmd

import (
	"bytes"
	"fmt"
	"os"
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
	Short: "{{ .APP }} version",
	Long:  `{{ .APP }} version`,
	Run: func(cmd *cobra.Command, args []string) {
        printVersion()
    },
}

func printVersion() {
	var versionBuffer bytes.Buffer

	if Version == "" {
		Version = "unknown"
	}
	versionBuffer.WriteString(fmt.Sprintf("{{ .APP }} version %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH))

	versionBuffer.WriteString(fmt.Sprintf("Go version %s\n", runtime.Version()))

	if Commit == "" {
		Commit = "unknown"
	}
	versionBuffer.WriteString(fmt.Sprintf("Git commit %s\n", Commit))

	if Date != "" {
		Date = cast.ToString(cast.ToTimeInDefaultLocation(Date, time.Local))
	} else {
		Date = "unknown"
	}
	versionBuffer.WriteString(fmt.Sprintf("Build date: %s\n", Date))

	fmt.Print(versionBuffer.String())
}

func init() {
	_ = os.Setenv("VERSION", Version)
	_ = os.Setenv("COMMIT", Commit)
	_ = os.Setenv("DATE", Date)

	rootCmd.AddCommand(versionCmd)
}
