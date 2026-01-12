/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package skills

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
)

var output string

// skillsCmd represents the skills command
var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: `Copy skills templates to .claude/skills directory`,
	Long:  `Copy skills templates from cmd/jzero/.template/skills to the project root's .claude/skills directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run()
	},
}

func GetCommand() *cobra.Command {
	skillsCmd.Flags().StringVarP(&output, "output", "o", ".claude/skills", "output directory")
	return skillsCmd
}

func Run() error {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Source template directory
	sourceDir := "skills"

	// Target directory in project root
	targetDir := filepath.Join(wd, output)

	// Check if source template exists
	entries := embeded.ReadTemplateDir(sourceDir)
	if entries == nil {
		return errors.New("skills templates not found")
	}

	// Copy template directory to target
	err = embeded.WriteTemplateDir(sourceDir, targetDir)
	if err != nil {
		return fmt.Errorf("failed to copy skills templates: %w", err)
	}

	fmt.Printf("Skills templates copied successfully to: %s\n", targetDir)
	return nil
}
