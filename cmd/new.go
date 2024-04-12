/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/jaronnie/jzero/embedx"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

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

func newProject(_ *cobra.Command, _ []string) error {
	// mkdir output
	err := os.MkdirAll(Dir, os.ModePerm)
	cobra.CheckErr(err)
	// go mod init
	_, err = Run(fmt.Sprintf("go mod init %s", Module), Dir)
	cobra.CheckErr(err)
	// touch main.go
	mainFile, err := ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embedx.ReadTemplateFile("jzero/main.go.tpl"))
	err = os.WriteFile(filepath.Join(Dir, "main.go"), mainFile, os.ModePerm)
	cobra.CheckErr(err)
	// mkdir cmd dir
	err = os.MkdirAll(filepath.Join(Dir, "cmd"), os.ModePerm)
	cobra.CheckErr(err)
	// touch cmd/root.go
	rootCmdFile, err := ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embedx.ReadTemplateFile("jzero/cmd/root.go.tpl"))
	err = os.WriteFile(filepath.Join(Dir, "cmd", "root.go"), rootCmdFile, os.ModePerm)
	cobra.CheckErr(err)
	// touch cmd/daemon.go
	daemonCmdFile, err := ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embedx.ReadTemplateFile("jzero/cmd/daemon.go.tpl"))
	err = os.WriteFile(filepath.Join(Dir, "cmd", "daemon.go"), daemonCmdFile, os.ModePerm)
	cobra.CheckErr(err)
	// mkdir daemon dir
	err = os.MkdirAll(filepath.Join(Dir, "daemon"), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/daemon.go
	daemonFile, err := ParseTemplate(map[string]interface{}{}, embedx.ReadTemplateFile("jzero/daemon/daemon.go.tpl"))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "daemon.go"), daemonFile, os.ModePerm)
	cobra.CheckErr(err)
	// mkdir api, proto dir
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "proto"), os.ModePerm)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "api"), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/api/{{.APP}}.api
	err = os.WriteFile(filepath.Join(Dir, "daemon", "api", APP+".api"), embedx.ReadTemplateFile("jzero/daemon/api/jzero.api.tpl"), os.ModePerm)
	cobra.CheckErr(err)
	// touch daemon/api/hello.api
	helloApiFile, err := ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embedx.ReadTemplateFile("jzero/daemon/api/hello.api.tpl"))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "api", "hello.api"), helloApiFile, os.ModePerm)
	cobra.CheckErr(err)

	// write proto dir
	err = embedx.WriteTemplateDir("jzero/daemon/proto", filepath.Join(Dir, "daemon", "proto"))
	cobra.CheckErr(err)

	// write config.toml
	configFile, err := ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embedx.ReadTemplateFile("jzero/config.toml.tpl"))
	err = os.WriteFile(filepath.Join(Dir, "config.toml"), configFile, os.ModePerm)
	cobra.CheckErr(err)

	// write .template
	//err = embedx.WriteTemplateDir("go-zero", filepath.Join(Dir, ".template", "go-zero"))
	//cobra.CheckErr(err)

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

var registerFuncMap = map[string]interface{}{}
