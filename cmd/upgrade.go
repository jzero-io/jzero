/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"github.com/jzero-io/jzero/config"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: `Upgrade the version of jzero tool.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch config.C.Upgrade.Channel {
		case "stable":
			return golang.Install("github.com/jzero-io/jzero@latest")
		case "main":
			return golang.Install("github.com/jzero-io/jzero@main")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringP("channel", "c", "main", "channel to upgrade jzero")
}
