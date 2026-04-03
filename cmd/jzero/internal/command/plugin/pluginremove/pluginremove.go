/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package pluginremove

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	Force bool
)

func Run(pluginName string) error {
	// Validate plugin name
	if !hasValidPrefix(pluginName, []string{"jzero"}) {
		return fmt.Errorf("plugin name must start with \"jzero-\" prefix. Got: %s", pluginName)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	pluginPath := filepath.Join(homeDir, ".jzero", "plugins", pluginName)

	// Check if plugin exists
	if _, err := os.Stat(pluginPath); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("plugin %s not found at %s", pluginName, pluginPath)
	}

	// Confirm removal
	if !Force {
		fmt.Fprintf(os.Stderr, "Are you sure you want to remove plugin %s? [y/N]: ", pluginName)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Fprint(os.Stderr, "Removal cancelled\n")
			return nil
		}
	}

	// Remove plugin
	if err := os.Remove(pluginPath); err != nil {
		return fmt.Errorf("failed to remove plugin: %w", err)
	}

	logx.Infof("Successfully removed plugin %s", pluginName)
	return nil
}

// hasValidPrefix checks if the plugin name has a valid prefix
func hasValidPrefix(name string, validPrefixes []string) bool {
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(name, prefix+"-") {
			return true
		}
	}
	return false
}
