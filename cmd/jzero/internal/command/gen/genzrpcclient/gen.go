package genzrpcclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/new"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/plugin"
)

type DirContext struct {
	ImportBase      string
	PbPackage       string
	OptionGoPackage string
	Resource        string
	Output          string
}

func (d DirContext) GetCall() generator.Dir {
	fileName := filepath.Join(d.Output, "typed", d.Resource)
	return generator.Dir{
		Filename: fileName,
		GetChildPackage: func(childPath string) (string, error) {
			return strings.ToLower(childPath), nil
		},
	}
}

func (d DirContext) GetEtc() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetInternal() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetConfig() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetLogic() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetServer() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetSvc() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetPb() generator.Dir {
	return generator.Dir{
		Package: d.packagePath(),
	}
}

func (d DirContext) packagePath() string {
	packagePath := filepath.ToSlash(fmt.Sprintf("%s/model%s/%s", d.ImportBase, d.Resource, strings.TrimPrefix(d.OptionGoPackage, "./")))
	return packagePath
}

func (d DirContext) GetProtoGo() generator.Dir {
	return generator.Dir{
		Filename: d.OptionGoPackage,
		Package:  d.packagePath(),
	}
}

func (d DirContext) GetMain() generator.Dir {
	panic("implement me")
}

func (d DirContext) GetServiceName() stringx.String {
	panic("implement me")
}

func (d DirContext) SetPbDir(pbDir, grpcDir string) {
	panic("implement me")
}

func Generate(genModule bool) (err error) {
	g := generator.NewGenerator(config.C.Gen.Style, false)

	var files []string

	switch {
	case len(config.C.Gen.Zrpcclient.Desc) > 0:
		for _, v := range config.C.Gen.Zrpcclient.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					files = append(files, v)
				}
			} else {
				specifiedProtoFiles, err := desc.GetProtoFilepath(v)
				if err != nil {
					return err
				}
				files = append(files, specifiedProtoFiles...)
			}
		}
	default:
		files, err = desc.GetProtoFilepath(config.C.ProtoDir())
		if err != nil {
			return err
		}
	}

	for _, v := range config.C.Gen.Zrpcclient.DescIgnore {
		if !osx.IsDir(v) {
			if filepath.Ext(v) == ".proto" {
				files = lo.Reject(files, func(item string, _ int) bool {
					return item == v
				})
			}
		} else {
			specifiedProtoFiles, err := desc.GetProtoFilepath(v)
			if err != nil {
				return err
			}
			for _, saf := range specifiedProtoFiles {
				files = lo.Reject(files, func(item string, _ int) bool {
					return item == saf
				})
			}
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 获取所有插件信息，用于后续判断文件来源
	plugins, _ := plugin.GetPlugins()

	var services []string
	for _, fp := range files {
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(fp, true)
		if err != nil {
			return err
		}
		dirContext := DirContext{
			ImportBase:      filepath.Join(config.C.Gen.Zrpcclient.GoModule),
			PbPackage:       parse.PbPackage,
			OptionGoPackage: parse.GoPackage,
			Output:          config.C.Gen.Zrpcclient.Output,
		}
		for _, service := range parse.Service {
			services = append(services, service.Name)
			_ = os.MkdirAll(filepath.Join(dirContext.GetCall().Filename, strings.ToLower(service.Name)), 0o755)
		}
		pbDir := filepath.Join(config.C.Gen.Zrpcclient.Output, "model")
		// gen pb model
		err = os.MkdirAll(pbDir, 0o755)
		if err != nil {
			return err
		}

		// 构建 protoc 命令，包含所有可能的导入路径
		var importPaths []string
		importPaths = append(importPaths, config.C.ProtoDir())
		importPaths = append(importPaths, filepath.Join(config.C.ProtoDir(), "third_party"))

		// 检查当前文件是否来自插件，如果是，添加对应的插件导入路径
		for _, p := range plugins {
			pluginProtoDir := filepath.Join(p.Path, "desc", "proto")
			if pathx.FileExists(pluginProtoDir) {
				importPaths = append(importPaths, pluginProtoDir)
				importPaths = append(importPaths, filepath.Join(pluginProtoDir, "third_party"))
			}
		}

		// 构建 -I 参数
		var protocolIncludePaths []string
		for _, path := range importPaths {
			protocolIncludePaths = append(protocolIncludePaths, fmt.Sprintf("-I%s", path))
		}

		protocCmd := fmt.Sprintf("protoc %s --go_out=%s --go-grpc_out=%s %s",
			strings.Join(protocolIncludePaths, " "), pbDir, pbDir, fp)
		resp, err := execx.Run(protocCmd, wd)
		if err != nil {
			return errors.Errorf("err: [%v], resp: [%s]", err, resp)
		}

		err = g.GenCall(dirContext, parse, &conf.Config{
			NamingFormat: config.C.Gen.Style,
		}, &generator.ZRpcContext{
			Multiple:    true,
			IsGenClient: true,
		})
		if err != nil {
			return err
		}
	}

	// 检查是否有插件
	hasPlugins := len(plugins) > 0
	for _, p := range plugins {
		if pathx.FileExists(filepath.Join(p.Path, "desc", "proto")) {
			pluginProtoFiles, err := desc.GetProtoFilepath(filepath.Join(p.Path, "desc", "proto"))
			if err == nil && len(pluginProtoFiles) > 0 {
				hasPlugins = true
				break
			}
		}
	}

	// gen clientset and options
	template, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "clientset.go.tpl")), map[string]any{
		"Module":     config.C.Gen.Zrpcclient.GoModule,
		"Package":    config.C.Gen.Zrpcclient.GoPackage,
		"Services":   services,
		"HasPlugins": hasPlugins,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "clientset.go.tpl"))))
	if err != nil {
		return err
	}

	formated, err := gosimports.Process("", template, nil)
	if err != nil {
		return errors.Errorf("format go file %s %s meet error: %v", filepath.Join(config.C.Gen.Zrpcclient.Output, "clientset.go"), template, err)
	}
	err = os.WriteFile(filepath.Join(config.C.Gen.Zrpcclient.Output, "clientset.go"), formated, 0o644)
	if err != nil {
		return err
	}

	// if set --module flag
	if genModule {
		data, err := new.NewTemplateData()
		if err != nil {
			return err
		}
		data["Module"] = config.C.Gen.Zrpcclient.GoModule
		if config.C.Gen.Zrpcclient.GoVersion != "" {
			data["GoVersion"] = config.C.Gen.Zrpcclient.GoVersion
		}
		template, err = templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl")), data, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl"))))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(config.C.Gen.Zrpcclient.Output, "go.mod"), template, 0o644)
		if err != nil {
			return err
		}
	}

	// 生成插件相关文件
	err = generatePluginFiles(plugins, config.C.Gen.Zrpcclient.GoModule, config.C.Gen.Zrpcclient.Output)
	if err != nil {
		return err
	}

	return nil
}

func generatePluginFiles(plugins []plugin.Plugin, goModule, output string) error {
	if len(plugins) == 0 {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	var pluginNames []string

	// 为每个插件生成完整的文件结构
	for _, p := range plugins {
		if !pathx.FileExists(filepath.Join(p.Path, "desc", "proto")) {
			continue
		}

		// 获取插件的 proto 文件
		pluginProtoFiles, err := desc.GetProtoFilepath(filepath.Join(p.Path, "desc", "proto"))
		if err != nil || len(pluginProtoFiles) == 0 {
			continue
		}

		var pluginServices []string

		// 1. 为插件生成 protobuf model 文件
		pluginModelDir := filepath.Join(output, "plugins", p.Name, "model")
		err = os.MkdirAll(pluginModelDir, 0o755)
		if err != nil {
			return err
		}

		// 解析插件的服务并生成 model
		for _, fp := range pluginProtoFiles {
			parser := rpcparser.NewDefaultProtoParser()
			parse, err := parser.Parse(fp, true)
			if err != nil {
				continue
			}

			// 收集服务名称
			for _, service := range parse.Service {
				pluginServices = append(pluginServices, service.Name)
			}

			// 为每个插件的 proto 文件生成 Go 代码
			// 构建 protoc 命令，包含所有可能的导入路径
			var importPaths []string
			importPaths = append(importPaths, config.C.ProtoDir())
			importPaths = append(importPaths, filepath.Join(config.C.ProtoDir(), "third_party"))
			importPaths = append(importPaths, filepath.Join(p.Path, "desc", "proto"))
			importPaths = append(importPaths, filepath.Join(p.Path, "desc", "proto", "third_party"))

			// 构建 -I 参数
			var protocolIncludePaths []string
			for _, path := range importPaths {
				protocolIncludePaths = append(protocolIncludePaths, fmt.Sprintf("-I%s", path))
			}

			protocCmd := fmt.Sprintf("protoc %s --go_out=%s --go-grpc_out=%s %s",
				strings.Join(protocolIncludePaths, " "), pluginModelDir, pluginModelDir, fp)
			resp, err := execx.Run(protocCmd, wd)
			if err != nil {
				return errors.Errorf("err: [%v], resp: [%s]", err, resp)
			}
		}

		if len(pluginServices) == 0 {
			continue
		}

		// 2. 为插件生成 typed 文件 (使用 go-zero 的生成器)
		g := generator.NewGenerator(config.C.Gen.Style, false)

		for _, fp := range pluginProtoFiles {
			parser := rpcparser.NewDefaultProtoParser()
			parse, err := parser.Parse(fp, true)
			if err != nil {
				continue
			}

			// 为每个服务创建目录（类似主服务的处理方式）
			pluginDirContext := DirContext{
				ImportBase:      filepath.Join(goModule, "plugins", p.Name),
				PbPackage:       parse.PbPackage,
				OptionGoPackage: parse.GoPackage,
				Output:          filepath.Join(output, "plugins", p.Name),
			}

			for _, service := range parse.Service {
				_ = os.MkdirAll(filepath.Join(pluginDirContext.GetCall().Filename, strings.ToLower(service.Name)), 0o755)
			}

			// 为整个 proto 文件生成客户端代码（类似主服务的处理方式）
			err = g.GenCall(pluginDirContext, parse, &conf.Config{
				NamingFormat: config.C.Gen.Style,
			}, &generator.ZRpcContext{
				Multiple:    true,
				IsGenClient: true,
			})
			if err != nil {
				return err
			}
		}

		pluginNames = append(pluginNames, p.Name)

		// 3. 生成单个插件的聚合文件
		pluginTemplate, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugin.go.tpl")), map[string]any{
			"Module":     goModule,
			"PluginName": p.Name,
			"Services":   pluginServices,
		}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugin.go.tpl"))))
		if err != nil {
			return err
		}

		formated, err := gosimports.Process("", pluginTemplate, nil)
		if err != nil {
			return errors.Errorf("format plugin go file meet error: %v", err)
		}

		pluginDir := filepath.Join(output, "plugins")
		err = os.MkdirAll(pluginDir, 0o755)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(pluginDir, p.Name+".go"), formated, 0o644)
		if err != nil {
			return err
		}
	}

	// 4. 生成主 plugins.go 文件
	pluginsTemplate, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugins.go.tpl")), map[string]any{
		"Module":      goModule,
		"PluginNames": pluginNames,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugins.go.tpl"))))
	if err != nil {
		return err
	}

	formated, err := gosimports.Process("", pluginsTemplate, nil)
	if err != nil {
		return errors.Errorf("format plugins go file meet error: %v", err)
	}

	pluginDir := filepath.Join(output, "plugins")
	err = os.MkdirAll(pluginDir, 0o755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(pluginDir, "plugins.go"), formated, 0o644)
	if err != nil {
		return err
	}

	return nil
}
