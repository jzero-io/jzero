/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package plugin

import (
	"fmt"
	"os/exec"
	"syscall"
)

// Handler is capable of parsing command line arguments
// and performing executable filename lookups to search
// for valid plugin files, and execute found plugins.
type Handler interface {
	// Lookup searches for a plugin executable with the given filename
	// and returns its full path if found, or a boolean false if not found.
	// Lookup will iterate over a list of given prefixes
	// in order to recognize valid plugin filenames.
	// The first filepath to match a prefix is returned.
	Lookup(filename string) (string, bool)

	// Execute receives an executable's filepath, a slice
	// of arguments, and a slice of environment variables
	// to relay to the executable.
	Execute(executablePath string, cmdArgs, environment []string) error
}

// DefaultHandler implements Handler
type DefaultHandler struct {
	ValidPrefixes []string
}

// NewDefaultHandler creates a new default plugin handler
func NewDefaultHandler(validPrefixes []string) *DefaultHandler {
	return &DefaultHandler{
		ValidPrefixes: validPrefixes,
	}
}

// Lookup implements Handler
func (h *DefaultHandler) Lookup(filename string) (string, bool) {
	// Search PATH for plugins with valid prefix
	for _, prefix := range h.ValidPrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-%s", prefix, filename))
		if err == nil && len(path) > 0 {
			return path, true
		}
	}

	return "", false
}

// Execute implements Handler
func (h *DefaultHandler) Execute(executablePath string, cmdArgs, environment []string) error {
	return syscall.Exec(executablePath, cmdArgs, environment)
}
