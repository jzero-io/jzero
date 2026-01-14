package skillsinit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
)

func Run() error {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Target directory in project root
	targetDir := filepath.Join(wd, config.C.Skills.Init.Output)

	// Copy template directory to target
	err = embeded.WriteTemplateDir("skills", targetDir)
	if err != nil {
		return fmt.Errorf("failed to initialized skills templates: %w", err)
	}

	if !config.C.Quiet {
		fmt.Printf("✓ Skills templates initialized successfully at: %s\n", config.C.Skills.Init.Output)
	}

	return nil
}
