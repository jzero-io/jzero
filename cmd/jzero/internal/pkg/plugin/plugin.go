package plugin

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	PluginPrefix = "jzero-"
)

func DiscoverPlugins() ([]*cobra.Command, error) {
	var commands []*cobra.Command

	path := os.Getenv("PATH")
	pathDirs := filepath.SplitList(path)

	pluginMap := make(map[string]bool)

	for _, dir := range pathDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			name := file.Name()
			if !strings.HasPrefix(name, PluginPrefix) {
				continue
			}

			if pluginMap[name] {
				continue
			}
			pluginMap[name] = true

			pluginName := strings.TrimPrefix(name, PluginPrefix)

			cmd := &cobra.Command{
				Use:   pluginName,
				Short: fmt.Sprintf("Execute %s plugin", name),
				Long:  fmt.Sprintf("Execute the %s plugin with the provided arguments", name),
				Args:  cobra.ArbitraryArgs,
				Run: func(cmd *cobra.Command, args []string) {
					execPlugin(name, args)
				},
				DisableFlagParsing: true,
			}

			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

func execPlugin(pluginName string, args []string) {
	cmd := exec.Command(pluginName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error executing plugin %s: %v\n", pluginName, err)
		os.Exit(1)
	}
}
