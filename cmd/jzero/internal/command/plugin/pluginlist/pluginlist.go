/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package pluginlist

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	NameOnly bool
)

func Run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logx.Errorf("Failed to get home directory: %v", err)
		return err
	}

	pluginDir := filepath.Join(homeDir, ".jzero", "plugins")
	pluginsFound := false

	// Check if plugin directory exists
	if _, err := os.Stat(pluginDir); errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "Plugin directory %s does not exist. No plugins installed.\n", pluginDir)
		return nil
	}

	// Collect plugins first
	var plugins []string
	err = filepath.WalkDir(pluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file has valid prefix
		filename := filepath.Base(path)
		if !hasValidPrefix(filename, []string{"jzero"}) {
			return nil
		}

		plugins = append(plugins, path)
		return nil
	})

	if err != nil {
		logx.Errorf("Error walking plugin directory: %v", err)
		return err
	}

	// Sort plugins
	sort.Strings(plugins)

	// Display plugins
	if len(plugins) > 0 {
		fmt.Fprint(os.Stderr, "The following jzero plugins are installed:\n\n")
		pluginsFound = true

		for _, pluginPath := range plugins {
			filename := filepath.Base(pluginPath)
			if NameOnly {
				fmt.Fprintf(os.Stdout, "%s\n", filename)
			} else {
				fmt.Fprintf(os.Stdout, "%s\n", pluginPath)
			}
		}
	}

	if !pluginsFound {
		fmt.Fprintf(os.Stderr, "No jzero plugins found in %s\n", pluginDir)
		fmt.Fprint(os.Stderr, "\nPlugins must be prefixed with \"jzero-\" to be recognized.\n")
	}

	return nil
}

// hasValidPrefix checks if the filepath has a valid prefix
func hasValidPrefix(filepath string, validPrefixes []string) bool {
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(filepath, prefix+"-") {
			return true
		}
	}
	return false
}
