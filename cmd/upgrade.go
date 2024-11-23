/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"

	"github.com/jzero-io/jzero/config"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: `Upgrade the version of jzero tool.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch config.C.Upgrade.Channel {
		case "stable":
			return golang.Install("github.com/jzero-io/jzero@latest")
		default:
			return golang.Install(fmt.Sprintf("github.com/jzero-io/jzero@%s", config.C.Upgrade.Channel))
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringP("channel", "c", "stable", "channel to upgrade jzero")
}
