/*
Copyright © 2024 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(worktab completion bash)

# To load completions for each session, execute once:
Linux:
  $ worktab completion bash > /etc/bash_completion.d/worktab
MacOS:
  $ worktab completion bash > /usr/local/etc/bash_completion.d/worktab

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ worktab completion zsh > "${fpath[1]}/_worktab"

# You will need to start a new shell for this setup to take effect.

Fish:

$ worktab completion fish | source

# To load completions for each session, execute once:
$ worktab completion fish > ~/.config/fish/completions/worktab.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
