/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package plugin

import (
	"os"
	"strings"
)

// HandlePluginCommand receives a pluginHandler and command-line arguments
// and attempts to find a plugin executable that satisfies the given arguments.
func HandlePluginCommand(pluginHandler Handler, cmdArgs []string) error {
	var remainingArgs []string // all "non-flag" arguments

	for idx := range cmdArgs {
		if strings.HasPrefix(cmdArgs[idx], "-") {
			break
		}
		remainingArgs = append(remainingArgs, strings.ReplaceAll(cmdArgs[idx], "-", "_"))
	}

	foundBinaryPath := ""

	// attempt to find binary, starting at longest possible name with given cmdArgs
	for len(remainingArgs) > 0 {
		path, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
		if !found {
			remainingArgs = remainingArgs[:len(remainingArgs)-1]
			continue
		}

		foundBinaryPath = path
		break
	}

	if len(foundBinaryPath) == 0 {
		return nil
	}

	// invoke cmd binary relaying the current environment and args given
	// remainingArgs will always have at least one element.
	// exec will make remainingArgs[0] the "binary name".
	if err := pluginHandler.Execute(foundBinaryPath, append([]string{foundBinaryPath}, cmdArgs[len(remainingArgs):]...), os.Environ()); err != nil {
		return err
	}

	return nil
}
