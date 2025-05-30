/*
Copyright © 2024 jaronnie jaron@jaronnie.com

*/

package completion

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(jzero completion bash)

# To load completions for each session, execute once:
Linux:
  $ jzero completion bash > /etc/bash_completion.d/jzero
MacOS:
  $ jzero completion bash > /usr/local/etc/bash_completion.d/jzero

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ jzero completion zsh > "${fpath[1]}/_jzero"

# You will need to start a new shell for this setup to take effect.

Fish:

$ jzero completion fish | source

# To load completions for each session, execute once:
$ jzero completion fish > ~/.config/fish/completions/jzero.fish

PowerShell:
$ jzero completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func GetCommand() *cobra.Command {
	return completionCmd
}
