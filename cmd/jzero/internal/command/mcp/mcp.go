/*
Copyright © 2025 jaronnie <jaron@jaronnie.com>
*/

package mcp

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jaronnie/genius"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "mcp server for jzero",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(cmd.Root())
	},
}

// mcpTestCmd represents the mcp test command
var mcpTestCmd = &cobra.Command{
	Use:   "test",
	Short: "mcp server test",
	RunE: func(cmd *cobra.Command, args []string) error {
		jzeroMcpCmd := exec.Command("jzero", "mcp")

		stdin, err := jzeroMcpCmd.StdinPipe()
		if err != nil {
			return err
		}

		stdout, err := jzeroMcpCmd.StdoutPipe()
		if err != nil {
			return err
		}

		stderr, err := jzeroMcpCmd.StderrPipe()
		if err != nil {
			return err
		}

		if err := jzeroMcpCmd.Start(); err != nil {
			return err
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Printf("%s %s\n", color.WithColor("[SERVER OUTPUT]", color.FgGreen), scanner.Text())
				fmt.Printf("%s Enter your method (press Enter to send):\n", color.WithColor("[INPUT]", color.FgYellow))
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Printf("%s %s\n", color.WithColor("[SERVER ERROR]", color.FgRed), scanner.Text())
				fmt.Printf("%s Enter your method (press Enter to send):\n", color.WithColor("[INPUT]", color.FgYellow))
			}
		}()

		fmt.Printf("%s Enter your method (press Enter to send):\n", color.WithColor("[INPUT]", color.FgYellow))
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			_, _ = stdin.Write([]byte(text + "\n"))

			g, err := genius.NewFromRawJSON(scanner.Bytes())
			if err == nil {
				if method, ok := g.Get("method").(string); ok && method != "" {
					if strings.HasPrefix(method, "notifications") {
						fmt.Printf("%s Enter your method (press Enter to send):\n", color.WithColor("[INPUT]", color.FgYellow))
					}
				}
			}
		}

		if err := jzeroMcpCmd.Wait(); err != nil {
			return err
		}

		return nil
	},
}

func GetCommand() *cobra.Command {
	mcpCmd.AddCommand(mcpTestCmd)
	return mcpCmd
}
