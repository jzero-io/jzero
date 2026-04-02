/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package plugin

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/plugin/pluginlist"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/plugin/pluginremove"
)

var (
	// ValidPluginFilenamePrefixes defines the allowed plugin prefix to search
	ValidPluginFilenamePrefixes = []string{"jzero"}
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage jzero plugins",
	Long: `Provides utilities for interacting with plugins.
Plugins provide extended functionality that is not part of the major command-line distribution.
Plugins must be prefixed with "jzero-" and are stored in ~/.jzero/plugins.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed jzero plugins",
	Long: `List all installed plugins from ~/.jzero/plugins directory.
Plugins must be prefixed with "jzero-" to be recognized.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return pluginlist.Run()
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:   "remove <plugin-name>",
	Short: "Remove an installed jzero plugin",
	Long: `Remove an installed jzero plugin from ~/.jzero/plugins.

Example:
  jzero plugin remove jzero-myplugin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pluginremove.Run(args[0])
	},
}

func GetCommand() *cobra.Command {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)

	pluginListCmd.Flags().BoolVar(&pluginlist.NameOnly, "name-only", false, "If true, display only the binary name of each plugin, rather than its full path")
	pluginRemoveCmd.Flags().BoolVar(&pluginremove.Force, "force", false, "Force removal without confirmation")

	return pluginCmd
}
