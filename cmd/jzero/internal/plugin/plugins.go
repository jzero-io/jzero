package plugin

import (
	"os"
	"path/filepath"

	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genrpc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
)

type Plugin struct {
	Path   string
	Module string
}

func GetPlugins() ([]Plugin, error) {
	var plugins []Plugin
	dir, err := os.ReadDir("plugins")
	if err != nil {
		return nil, err
	}
	for _, p := range dir {
		if p.IsDir() {
			plugins = append(plugins, Plugin{
				Path: filepath.ToSlash(filepath.Join("plugins", p.Name())),
			})
		} else if p.Type() == os.ModeSymlink {
			plugins = append(plugins, Plugin{
				Path: filepath.ToSlash(filepath.Join("plugins", p.Name())),
			})
		}
	}
	return plugins, nil
}

func GetProjectType() (string, error) {
	// 判断 core 项目类型 api/rpc
	var projectType string
	if _, err := os.Stat(filepath.Join("desc", "api")); err == nil {
		// api 项目
		projectType = "api"
	}
	if _, err := os.Stat(filepath.Join("desc", "proto")); err == nil {
		// rpc 项目
		projectType = "rpc"

		// 获取全量 proto 文件
		protoFiles, err := jzerodesc.GetProtoFilepath(config.C.ProtoDir())
		if err != nil {
			return "", err
		}

		for _, v := range protoFiles {
			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			var parse rpcparser.Proto
			parse, err = protoParser.Parse(v, true)
			if err != nil {
				return "", err
			}
			if genrpc.IsNeedGenProtoDescriptor(parse) {
				projectType = "gateway"
				break
			}
		}
	}

	return projectType, nil
}
