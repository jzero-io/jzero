/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package skills

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/skills/skillsinit"
)

// skillsCmd represents the skills command
var skillsCmd = &cobra.Command{
	Use:  "skills",
	Long: `Manage skills for AI agents.`,
}

// initCmd represents the skills init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: `Initialize skills templates in .claude/skills directory`,
	Long:  `Initialize skills templates from cmd/jzero/.template/skills to the project root's .claude/skills directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return skillsinit.Run()
	},
}

func GetCommand() *cobra.Command {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}
	initCmd.Flags().StringP("output", "o", filepath.Join(homeDir, ".claude", "skills"), "output directory")

	// Add init subcommand
	skillsCmd.AddCommand(initCmd)

	return skillsCmd
}
