/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"bytes"
	"fmt"
	sprig "github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

/*
tpl
*/
var mainTpl = `package main

import (
	"{{ .Module }}/cmd"
)

func main() {
	cmd.Execute()
}

`

var rootCmdTpl = `package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .APP }}",
	Short: "{{ .APP }} root",
	Long:  "{{ .APP }} framework.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .APP }}/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(os.Args) >= 1 && os.Args[1] != daemonCmd.Name() {
		return
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".{{ .APP }}" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".{{ .APP }}"))
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		cfgFile = viper.ConfigFileUsed()
	} else {
		cobra.CheckErr(err)
	}
}
`

var daemonCmdTpl = `package cmd

import (
	"github.com/spf13/cobra"

	"{{ .Module }}/daemon"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "{{ .APP }} daemon",
	Long:  "{{ .APP }} daemon",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start(cfgFile)
		select {}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

`

var daemonTpl = `package daemon

import (
	"fmt"
)

func Start(cfgFile string) {
	go func() {
		fmt.Println("start")
	}()
}
`

var apiTpl = `syntax = "v1"
import "hello.api"

`

var helloApiTpl = `syntax = "v1"

type pathRequest struct {
    Name string ` + "`path:\"name\"`" + `
}

type paramRequest struct {
    Name string ` + "`form:\"name\"`" + `
}

type postRequest struct {
    Name string ` + "`json:\"name\"`" + `
}

type response {
    Message string
}

@server(
    prefix: /api/v1
    group: hello
)
service daemon {
    @handler HelloPathHandler
    get /hello/:name (pathRequest) returns (response)

    @handler HelloParamHandler
    get /hello (paramRequest) returns (response)

    @handler HelloPostHandler
    post /hello (postRequest) returns (response)
}
`

var credentialProtoTpl = `
`

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "jzero new project",
	Long:  `jzero new project`,
	RunE:  newProject,
}

var (
	Module string
	Dir    string
	APP    string
)

const (
	// OsWindows represents os windows
	OsWindows = "windows"
	// OsMac represents os mac
	OsMac = "darwin"
	// OsLinux represents os linux
	OsLinux = "linux"
)

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&Module, "module", "m", "", "set go module")
	newCmd.Flags().StringVarP(&Dir, "dir", "d", "", "set output dir")
	newCmd.Flags().StringVarP(&APP, "app", "", "", "set app name")
}

func newProject(cmd *cobra.Command, args []string) error {
	// mkdir output
	err := os.MkdirAll(Dir, os.ModePerm)
	cobra.CheckErr(err)
	// go mod init
	_, err = Run(fmt.Sprintf("go mod init %s", Module), Dir)
	cobra.CheckErr(err)
	// touch main.go
	mainFile, err := ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, []byte(mainTpl))
	err = os.WriteFile(filepath.Join(Dir, "main.go"), mainFile, os.ModePerm)
	cobra.CheckErr(err)
	// mkdir cmd dir
	err = os.MkdirAll(filepath.Join(Dir, "cmd"), os.ModePerm)
	cobra.CheckErr(err)
	// touch cmd/root.go
	rootCmdFile, err := ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, []byte(rootCmdTpl))
	err = os.WriteFile(filepath.Join(Dir, "cmd", "root.go"), rootCmdFile, os.ModePerm)
	cobra.CheckErr(err)
	// touch cmd/daemon.go
	daemonCmdFile, err := ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, []byte(daemonCmdTpl))
	err = os.WriteFile(filepath.Join(Dir, "cmd", "daemon.go"), daemonCmdFile, os.ModePerm)
	cobra.CheckErr(err)
	// mkdir daemon dir
	err = os.MkdirAll(filepath.Join(Dir, "daemon"), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/daemon.go
	daemonFile, err := ParseTemplate(map[string]interface{}{}, []byte(daemonTpl))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "daemon.go"), daemonFile, os.ModePerm)
	cobra.CheckErr(err)
	// mkdir api, proto dir
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "proto"), os.ModePerm)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "api"), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/api/{{.APP}}d.api
	err = os.WriteFile(filepath.Join(Dir, "daemon", "api", APP+"d.api"), []byte(apiTpl), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/api/hello.api
	err = os.WriteFile(filepath.Join(Dir, "daemon", "api", "hello.api"), []byte(helloApiTpl), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/proto/credential.proto
	// touch daemon/proto/machine.proto

	// run go mod tidy
	_, err = Run("go mod tidy", Dir)
	cobra.CheckErr(err)
	return nil
}

// ParseTemplate template
func ParseTemplate(data interface{}, tplT []byte) ([]byte, error) {
	t := template.Must(template.New("production").Funcs(sprig.TxtFuncMap()).Funcs(RegisterTxtFuncMap()).Parse(string(tplT)))

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Run provides the execution of shell scripts in golang,
// which can support macOS, Windows, and Linux operating systems.
// Other operating systems are currently not supported
func Run(arg, dir string, in ...*bytes.Buffer) (string, error) {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case OsMac, OsLinux:
		cmd = exec.Command("sh", "-c", arg)
	case OsWindows:
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return "", fmt.Errorf("unexpected os: %v", goos)
	}
	if len(dir) > 0 {
		cmd.Dir = dir
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if len(in) > 0 {
		cmd.Stdin = in[0]
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSuffix(stderr.String(), "\n"))
		}
		return "", err
	}

	return strings.TrimSuffix(stdout.String(), "\n"), nil
}

func RegisterTxtFuncMap() template.FuncMap {
	return RegisterFuncMap()
}

func RegisterFuncMap() template.FuncMap {
	gfm := make(map[string]interface{}, len(registerFuncMap))
	for k, v := range registerFuncMap {
		gfm[k] = v
	}
	return gfm
}

var registerFuncMap = map[string]interface{}{
	"firstUpper": FirstUpper,
	"firstLower": FirstLower,
}

func FirstUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

func FirstLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}
