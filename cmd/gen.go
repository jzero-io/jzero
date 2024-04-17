/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jaronnie/genius"
	"github.com/spf13/cast"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"

	"github.com/jaronnie/jzero/embeded"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "jzero gen code",
	Long:  `jzero gen code`,
	RunE:  gen,
}

type ImportLines []string

type RegisterLines []string

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

func gen(_ *cobra.Command, _ []string) error {
	wd, err := os.Getwd()
	cobra.CheckErr(err)

	moduleStruct, err := GetGoMod(wd)
	cobra.CheckErr(err)
	Module = moduleStruct.Path

	// read proto dir
	ds, err := os.ReadDir(filepath.Join(wd, "daemon", "proto"))
	cobra.CheckErr(err)

	var protosets []string
	var serverImports ImportLines
	var pbImports ImportLines
	var registerServers RegisterLines

	for _, v := range ds {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "proto") {
			command := fmt.Sprintf("goctl rpc protoc daemon/proto/%s  -I./daemon/proto --go_out=./daemon --go-grpc_out=./daemon  --zrpc_out=./daemon --client=false --home %s -m", v.Name(), filepath.Join(wd, ".template", "go-zero"))
			_, err = Run(command, wd)
			cobra.CheckErr(err)

			fileBase := v.Name()[0 : len(v.Name())-len(path.Ext(v.Name()))]
			_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", fileBase)))

			// # gen proto descriptor
			//protoc --include_imports -I./daemon/proto --descriptor_set_out=.protosets/xx.pb daemon/proto/xx.proto
			_ = os.MkdirAll(filepath.Join(wd, ".protosets"), 0o755)
			protocCommand := fmt.Sprintf("protoc --include_imports -I./daemon/proto --descriptor_set_out=.protosets/%s.pb daemon/proto/%s.proto", fileBase, fileBase)
			_, err = Run(protocCommand, wd)
			cobra.CheckErr(err)

			protosets = append(protosets, filepath.Join(".protosets", fmt.Sprintf("%s.pb", fileBase)))

			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			parse, err := protoParser.Parse(filepath.Join(wd, "daemon", "proto", v.Name()), true)
			cobra.CheckErr(err)
			for _, s := range parse.Service {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/daemon/internal/server/%s"`, s.Name, Module, s.Name))
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), FirstUpper(s.Name), s.Name, FirstUpper(s.Name)))
			}
			pbImports = append(pbImports, fmt.Sprintf(`"%s/daemon/%s"`, Module, strings.TrimPrefix(parse.GoPackage, "./")))
		}
	}

	// read api file
	configBytes, err := os.ReadFile(filepath.Join(wd, "config.toml"))
	cobra.CheckErr(err)

	// 修改 config.toml protosets 内容
	g, err := genius.NewFromToml(configBytes)
	cobra.CheckErr(err)
	err = g.Set("Gateway.Upstreams.0.ProtoSets", protosets)
	cobra.CheckErr(err)
	toml, err := g.EncodeToToml()
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(wd, "config.toml"), toml, 0644)
	cobra.CheckErr(err)

	// 修改 daemon/zrpc.go
	zrpcFile, err := ParseTemplate(map[string]interface{}{
		"Module":          Module,
		"APP":             APP,
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "zrpc.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "zrpc.go"), zrpcFile, 0o644)
	cobra.CheckErr(err)

	// 生成 api 代码
	command := fmt.Sprintf("goctl api go --api daemon/api/%s.api --dir ./daemon --home %s", cast.ToString(g.Get("APP")), filepath.Join(wd, ".template", "go-zero"))
	_, err = Run(command, wd)
	cobra.CheckErr(err)

	// 删除无用文件夹
	_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", cast.ToString(g.Get("APP")))))
	_ = os.RemoveAll(filepath.Join(wd, "daemon", "etc"))

	return nil
}

// GetGoMod is used to determine whether workDir is a go module project through command `go list -json -m`
func GetGoMod(workDir string) (*ModuleStruct, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	data, err := execx.Run("go list -json -m", workDir)
	if err != nil {
		return nil, nil
	}

	var m ModuleStruct
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// ModuleStruct contains the relative data of go module,
// which is the result of the command go list
type ModuleStruct struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

func FirstUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

func init() {
	rootCmd.AddCommand(genCmd)
}
