package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jaronnie/genius"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"

	"github.com/jaronnie/jzero/daemon/pkg/mod"
	"github.com/jaronnie/jzero/daemon/pkg/stringx"
	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"
)

var (
	WorkingDir string
)

type (
	ImportLines   []string
	RegisterLines []string
)

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

func Gen(c *cobra.Command, _ []string) error {
	// change dir
	if WorkingDir != "" {
		err := os.Chdir(WorkingDir)
		cobra.CheckErr(err)
	}

	wd, err := os.Getwd()
	cobra.CheckErr(err)

	moduleStruct, err := mod.GetGoMod(wd)
	cobra.CheckErr(err)

	// read proto dir
	ds, err := os.ReadDir(filepath.Join(wd, "daemon", "desc", "proto"))
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
			command := fmt.Sprintf("goctl rpc protoc daemon/desc/proto/%s  -I./daemon/desc/proto --go_out=./daemon --go-grpc_out=./daemon  --zrpc_out=./daemon --client=false --home %s -m", v.Name(), filepath.Join(wd, ".template", "go-zero"))
			_, err = execx.Run(command, wd)
			cobra.CheckErr(err)

			fileBase := v.Name()[0 : len(v.Name())-len(path.Ext(v.Name()))]
			_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", fileBase)))

			// # gen proto descriptor
			_ = os.MkdirAll(filepath.Join(wd, ".protosets"), 0o755)
			protocCommand := fmt.Sprintf("protoc --include_imports -I./daemon/desc/proto --descriptor_set_out=.protosets/%s.pb daemon/desc/proto/%s.proto", fileBase, fileBase)
			_, err = execx.Run(protocCommand, wd)
			cobra.CheckErr(err)

			protosets = append(protosets, filepath.Join(".protosets", fmt.Sprintf("%s.pb", fileBase)))

			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			parse, err := protoParser.Parse(filepath.Join(wd, "daemon", "desc", "proto", v.Name()), true)
			cobra.CheckErr(err)
			for _, s := range parse.Service {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/daemon/internal/server/%s"`, s.Name, moduleStruct.Path, s.Name))
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), s.Name, stringx.FirstUpper(s.Name)))
			}
			pbImports = append(pbImports, fmt.Sprintf(`"%s/daemon/%s"`, moduleStruct.Path, strings.TrimPrefix(parse.GoPackage, "./")))
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
	zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module":          moduleStruct.Path,
		"APP":             cast.ToString(g.Get("APP")),
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "zrpc.go.tpl")))
	err = os.WriteFile(filepath.Join(wd, "daemon", "zrpc.go"), zrpcFile, 0o644)
	cobra.CheckErr(err)

	// 生成 api 代码
	command := fmt.Sprintf("goctl api go --api daemon/desc/api/%s.api --dir ./daemon --home %s", cast.ToString(g.Get("APP")), filepath.Join(wd, ".template", "go-zero"))
	_, err = execx.Run(command, wd)
	cobra.CheckErr(err)

	// 删除无用文件夹
	_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", cast.ToString(g.Get("APP")))))
	_ = os.RemoveAll(filepath.Join(wd, "daemon", "etc"))

	return nil
}
