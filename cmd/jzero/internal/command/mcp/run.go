package mcp

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mcp"
)

func Run(rootCmd *cobra.Command) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	commands := rootCmd.Commands()
	for _, cmd := range commands {
		if cmd.Use == "mcp" {
			rootCmd.RemoveCommand(cmd)
		}
	}
	mcpServer := mcp.NewCobraMCPServer(rootCmd)
	go func() {
		if err := mcpServer.ServeStdio(); err != nil {
			fmt.Printf("MCP server error: %v\n", err)
			os.Exit(1)
		}
	}()
	<-quit
	return nil
}
