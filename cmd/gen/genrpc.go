package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jaronnie/genius"
	"github.com/jzero-io/jzero/app/pkg/stringx"
	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/jzero-io/jzero/embeded"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type JzeroRpc struct {
	Wd     string
	Module string
}

func (jr *JzeroRpc) Gen() error {
	protoDir, err := GetProtoDir(jr.Wd)
	if err != nil {
		return err
	}

	// get configType
	configType, err := stringx.GetConfigType(jr.Wd)
	if err != nil {
		return err
	}

	configBytes, err := os.ReadFile(filepath.Join(jr.Wd, "config."+configType))
	if err != nil {
		return err
	}

	g, err := genius.NewFromType(configBytes, configType)
	if err != nil {
		return err
	}

	var protosets []string
	var serverImports ImportLines
	var pbImports ImportLines
	var registerServers RegisterLines

	for _, v := range protoDir {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "proto") {
			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			parse, err := protoParser.Parse(filepath.Join(jr.Wd, "app", "desc", "proto", v.Name()), true)
			if err != nil {
				return err
			}

			fmt.Printf("%s to generate proto code. \n%s proto file %s\n", color.WithColor("Start", color.FgGreen), color.WithColor("Using", color.FgGreen), filepath.Join(jr.Wd, "app", "desc", "proto", v.Name()))
			command := fmt.Sprintf("goctl rpc protoc app/desc/proto/%s  -I./app/desc/proto --go_out=./app/internal --go-grpc_out=./app/internal --zrpc_out=./app --client=false --home %s -m", v.Name(), filepath.Join(embeded.Home, "go-zero"))
			_, err = execx.Run(command, jr.Wd)
			if err != nil {
				return err
			}
			fmt.Println(color.WithColor("Done", color.FgGreen))

			fileBase := v.Name()[0 : len(v.Name())-len(path.Ext(v.Name()))]
			rmf := strings.ReplaceAll(strings.ToLower(fileBase), "-", "")
			rmf = strings.ReplaceAll(rmf, "_", "")
			_ = os.Remove(filepath.Join(jr.Wd, "app", fmt.Sprintf("%s.go", rmf)))

			// # gen proto descriptor
			if isNeedGenProtoDescriptor(parse) {
				_ = os.MkdirAll(filepath.Join(jr.Wd, ".protosets"), 0o755)
				protocCommand := fmt.Sprintf("protoc --include_imports -I./app/desc/proto --descriptor_set_out=.protosets/%s.pb app/desc/proto/%s.proto", fileBase, fileBase)
				_, err = execx.Run(protocCommand, jr.Wd)
				if err != nil {
					return err
				}
				protosets = append(protosets, filepath.Join(".protosets", fmt.Sprintf("%s.pb", fileBase)))
			}

			for _, s := range parse.Service {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/app/internal/server/%s"`, s.Name, jr.Module, s.Name))
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), s.Name, stringx.FirstUpper(s.Name)))
			}
			pbImports = append(pbImports, fmt.Sprintf(`"%s/app/internal/%s"`, jr.Module, strings.TrimPrefix(parse.GoPackage, "./")))
		}
	}

	// 生成 app/zrpc.go
	if pathx.FileExists(filepath.Join(jr.Wd, "app", "zrpc.go")) {
		fmt.Printf("%s to generate app/zrpc.go\n", color.WithColor("Start", color.FgGreen))
		zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
			"Module":          jr.Module,
			"APP":             cast.ToString(g.Get("APP")),
			"ServerImports":   serverImports,
			"PbImports":       pbImports,
			"RegisterServers": registerServers,
		}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "zrpc.go.tpl")))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(jr.Wd, "app", "zrpc.go"), zrpcFile, 0o644)
		if err != nil {
			return err
		}
		fmt.Printf("%s", color.WithColor("Done\n", color.FgGreen))

		if g.Get("Gateway") != nil {
			// 修改 config.toml protosets 内容
			// 检测是否需要修改 config.toml. 以及让用户选择是否自动更新文件
			existProtosets := g.Get("Gateway.Upstreams.0.ProtoSets")
			if len(lo.Intersect(cast.ToStringSlice(existProtosets), protosets)) != len(protosets) {
				var in string
				fmt.Printf("检测到 config.%s 中 Gateway.Upstreams.0.ProtoSets 配置需要更新. 是否自动更新 y/n. 更新需谨慎, 会将注释删掉\n", configType)
				_, _ = fmt.Scanln(&in)
				switch {
				case strings.EqualFold(in, "y"):
					fmt.Printf("%s to update config.%s\n", color.WithColor("Start", color.FgGreen), configType)
					err = g.Set("Gateway.Upstreams.0.ProtoSets", protosets)
					if err != nil {
						return err
					}
					configBytes, err := g.EncodeToType(configType)
					if err != nil {
						return err
					}
					err = os.WriteFile(filepath.Join(jr.Wd, "config."+configType), configBytes, 0o644)
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
				case strings.EqualFold(in, "n"):
					fmt.Printf("请手动更新 Gateway.Upstreams.0.ProtoSets 配置\n配置该值为: \n%s\n",
						color.WithColor(fmt.Sprintf("%v", protosets), color.FgGreen))
				}
			}
		}
	}
	return nil
}

func isNeedGenProtoDescriptor(proto rpcparser.Proto) bool {
	for _, ps := range proto.Service {
		for _, rpc := range ps.RPC {
			for _, option := range rpc.Options {
				if option.Name == "(google.api.http)" {
					return true
				}
			}
		}
	}
	return false
}
