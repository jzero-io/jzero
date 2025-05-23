package mcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type MCPServer struct {
	rootCmd *cobra.Command
	server  *server.MCPServer
}

func NewMCPServer(rootCmd *cobra.Command) *MCPServer {
	s := &MCPServer{
		rootCmd: rootCmd,
		server: server.NewMCPServer(rootCmd.Short, "1.0.0",
			server.WithLogging(),
			server.WithRecovery(),
			server.WithResourceCapabilities(true, true),
		),
	}

	leaves := getLeafCommands(rootCmd)

	for _, leaf := range leaves {
		fullPath := getFullCommandPath(leaf)
		toolName := strings.Join(fullPath, "__")
		toolDesc := leaf.Short
		if toolDesc == "" {
			toolDesc = leaf.Long
		}
		var toolOptions []mcp.ToolOption
		toolOptions = append(toolOptions, mcp.WithDescription(toolDesc))

		flags := getAllFlagDefs(leaf)
		for _, f := range flags {
			var propertiesOptions []mcp.PropertyOption
			if f.DefValue == "" {
				propertiesOptions = append(propertiesOptions, mcp.Required())
			}
			propertiesOptions = append(propertiesOptions, mcp.Description(f.Usage))
			switch f.Value.Type() {
			case "string":
				propertiesOptions = append(propertiesOptions, mcp.DefaultString(f.DefValue))
				toolOptions = append(toolOptions, mcp.WithString(f.Name, propertiesOptions...))
			case "int":
				defaultValue, _ := strconv.ParseFloat(f.DefValue, 64)
				propertiesOptions = append(propertiesOptions, mcp.DefaultNumber(defaultValue))
				toolOptions = append(toolOptions, mcp.WithNumber(f.Name, propertiesOptions...))
			case "bool":
				defaultValue, _ := strconv.ParseBool(f.DefValue)
				propertiesOptions = append(propertiesOptions, mcp.DefaultBool(defaultValue))
				toolOptions = append(toolOptions, mcp.WithBoolean(f.Name, propertiesOptions...))
			case "float32", "float64":
				defaultValue, _ := strconv.ParseFloat(f.DefValue, 64)
				propertiesOptions = append(propertiesOptions, mcp.DefaultNumber(defaultValue))
				toolOptions = append(toolOptions, mcp.WithNumber(f.Name, propertiesOptions...))
			default:
				toolOptions = append(toolOptions, mcp.WithString(f.Name, propertiesOptions...))
			}
		}
		tool := mcp.NewTool(toolName, toolOptions...)

		s.server.AddTool(tool, s.handleToolCall(leaf))
	}

	return s
}

func (s *MCPServer) ServeStdio() error {
	return server.ServeStdio(s.server)
}

func (s *MCPServer) handleToolCall(cmd *cobra.Command) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		var fullArgs []string
		fullArgs = append(fullArgs, getFullCommandPath(cmd)...)

		for key, val := range request.GetArguments() {
			if key == "args" {
				continue
			}
			fullArgs = append(fullArgs, "--"+key)
			fullArgs = append(fullArgs, fmt.Sprintf("%v", val))
		}

		if args, ok := request.GetArguments()["args"].([]any); ok {
			for _, arg := range args {
				fullArgs = append(fullArgs, fmt.Sprintf("%v", arg))
			}
		}

		s.rootCmd.SetArgs(fullArgs)

		err := s.rootCmd.Execute()

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = originalStdout
		if err != nil {
			return mcp.NewToolResultText(err.Error()), nil
		}
		capturedText := buf.String()
		return mcp.NewToolResultText(capturedText), nil
	}
}

func getLeafCommands(cmd *cobra.Command) []*cobra.Command {
	var leaves []*cobra.Command
	if len(cmd.Commands()) == 0 {
		leaves = append(leaves, cmd)
	} else {
		leaves = append(leaves, cmd)
		for _, sub := range cmd.Commands() {
			leaves = append(leaves, getLeafCommands(sub)...)
		}
	}
	return leaves
}

func getFullCommandPath(cmd *cobra.Command) []string {
	if cmd.Parent() == nil {
		// ignore the root command
		return []string{}
	}
	parentPath := getFullCommandPath(cmd.Parent())
	return append(parentPath, cmd.Name())
}

func getAllFlagDefs(cmd *cobra.Command) map[string]*pflag.Flag {
	flags := make(map[string]*pflag.Flag)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = f
	})
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = f
	})
	return flags
}
