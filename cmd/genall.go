/*
Copyright © 2025 jaron <jaron@jaronnie.com>
*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/internal/genall"
)

// genallCmd represents the genall command
var genallCmd = &cobra.Command{
	Use:   "genall",
	Short: "jzero gen all codes which contains api/types/handler/model/logic file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return genall.Run()
	},
}

func init() {
	rootCmd.AddCommand(genallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
