package genzrpcclient

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
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
	g := generator.NewGenerator(config.C.Style, false)

	var files []string

	switch {
	case len(config.C.Gen.Zrpcclient.Desc) > 0:
		for _, v := range config.C.Gen.Zrpcclient.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					files = append(files, v)
				}
			} else {
				specifiedProtoFiles, err := desc.FindRpcServiceProtoFiles(v)
				if err != nil {
					return err
				}
				files = append(files, specifiedProtoFiles...)
			}
		}
	default:
		files, err = desc.FindRpcServiceProtoFiles(config.C.ProtoDir())
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
			specifiedProtoFiles, err := desc.FindRpcServiceProtoFiles(v)
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

	plugins, _ := plugin.GetPlugins()

	excludeThirdPartyProtoFiles, err := desc.FindExcludeThirdPartyProtoFiles(config.C.ProtoDir())
	if err != nil {
		return err
	}
	logx.Debugf("excludeThirdPartyProtoFiles: %v", excludeThirdPartyProtoFiles)

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

		var importPaths []string
		importPaths = append(importPaths, config.C.ProtoDir())
		importPaths = append(importPaths, filepath.Join(config.C.ProtoDir(), "third_party"))

		var protoParser protoparse.Parser
		protoParser.InferImportPaths = false

		protoDir := filepath.Join("desc", "proto")
		thirdPartyProtoDir := filepath.Join("desc", "proto", "third_party")
		protoParser.ImportPaths = []string{protoDir, thirdPartyProtoDir}
		for _, v := range config.C.Gen.Zrpcclient.ProtoInclude {
			protoParser.ImportPaths = append(protoParser.ImportPaths, v)
		}
		protoParser.IncludeSourceCodeInfo = true

		rel, err := filepath.Rel(config.C.ProtoDir(), fp)
		if err != nil {
			return err
		}

		fds, err := protoParser.ParseFiles(rel)
		if err != nil {
			return err
		}

		if len(fds) == 0 {
			continue
		}

		goPackage := fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

		getMod, err := mod.GetGoMod(config.C.Wd())
		if err != nil {
			return err
		}

		module := config.C.Gen.Zrpcclient.GoModule
		if !genModule {
			if config.C.Gen.Zrpcclient.Output != "." {
				module = getMod.Path
			}
		}

		protocCmd := fmt.Sprintf("protoc %s -I%s -I%s --go_out=%s --go-grpc_out=%s",
			fp,
			config.C.ProtoDir(),
			filepath.Join(config.C.ProtoDir(), "third_party"),
			func() string {
				if !genModule {
					return "."
				}
				return filepath.Join(config.C.Gen.Zrpcclient.Output)
			}(),
			func() string {
				if !genModule {
					return "."
				}
				return filepath.Join(config.C.Gen.Zrpcclient.Output)
			}(),
		)

		for _, exp := range excludeThirdPartyProtoFiles {
			rel, err = filepath.Rel(config.C.ProtoDir(), exp)
			if err != nil {
				return err
			}

			fds, err = protoParser.ParseFiles(rel)
			if err != nil {
				return err
			}

			if len(fds) == 0 {
				continue
			}

			goPackage = fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

			protocCmd += fmt.Sprintf(" --go_opt=module=%s --go_opt=M%s=%s --go-grpc_opt=module=%s --go-grpc_opt=M%s=%s", module, rel, func() string {
				if strings.HasPrefix(goPackage, module) {
					return goPackage
				}

				if genModule {
					return filepath.ToSlash(filepath.Join(module, "model", goPackage))
				}
				return filepath.ToSlash(filepath.Join(module, config.C.Gen.Zrpcclient.Output, "model", goPackage))
			}(), module, rel, func() string {
				if strings.HasPrefix(goPackage, module) {
					return goPackage
				}

				if genModule {
					return filepath.ToSlash(filepath.Join(module, "model", goPackage))
				}
				return filepath.ToSlash(filepath.Join(module, config.C.Gen.Zrpcclient.Output, "model", goPackage))
			}())
		}

		if len(config.C.Gen.Zrpcclient.ProtoInclude) > 0 {
			protocCmd += fmt.Sprintf(" -I%s ", strings.Join(config.C.Gen.Zrpcclient.ProtoInclude, " -I"))
		}

		logx.Debugf(protocCmd)
		resp, err := execx.Run(protocCmd, wd)
		if err != nil {
			return errors.Errorf("err: [%v], resp: [%s]", err, resp)
		}

		err = g.GenCall(dirContext, parse, &conf.Config{
			NamingFormat: config.C.Style,
		}, &generator.ZRpcContext{
			Multiple:    true,
			IsGenClient: true,
		})
		if err != nil {
			return err
		}
	}

	hasPlugins := len(plugins) > 0
	for _, p := range plugins {
		if pathx.FileExists(filepath.Join(p.Path, "desc", "proto")) {
			pluginProtoFiles, err := desc.FindRpcServiceProtoFiles(filepath.Join(p.Path, "desc", "proto"))
			if err == nil && len(pluginProtoFiles) > 0 {
				hasPlugins = true
				break
			}
		}
	}

	// gen clientset
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
		goVersion, err := mod.GetGoVersion()
		if err != nil {
			return err
		}
		templateData := map[string]any{
			"GoVersion": goVersion,
			"GoArch":    runtime.GOARCH,
		}
		templateData["Module"] = config.C.Gen.Zrpcclient.GoModule
		if config.C.Gen.Zrpcclient.GoVersion != "" {
			templateData["GoVersion"] = config.C.Gen.Zrpcclient.GoVersion
		}
		template, err = templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl")), templateData, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl"))))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(config.C.Gen.Zrpcclient.Output, "go.mod"), template, 0o644)
		if err != nil {
			return err
		}
	}

	err = generatePluginFiles(plugins)
	if err != nil {
		return err
	}

	err = genNoRpcServiceExcludeThirdPartyProto(genModule, config.C.Gen.Zrpcclient.GoModule)
	if err != nil {
		return err
	}

	return nil
}

func generatePluginFiles(plugins []plugin.Plugin) error {
	if len(plugins) == 0 {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	var pluginNames []string

	for _, p := range plugins {
		pluginProtoDir := filepath.Join(p.Path, "desc", "proto")
		pluginThirdPartyProtoDir := filepath.Join(p.Path, "desc", "proto", "third_party")

		if !pathx.FileExists(pluginProtoDir) {
			continue
		}

		pluginProtoFiles, err := desc.FindRpcServiceProtoFiles(pluginProtoDir)
		if err != nil || len(pluginProtoFiles) == 0 {
			continue
		}

		var pluginServices []string

		pluginModelDir := filepath.Join(config.C.Gen.Zrpcclient.Output, "plugins", p.Name, "model")
		err = os.MkdirAll(pluginModelDir, 0o755)
		if err != nil {
			return err
		}

		excludeThirdPartyProtoFiles, err := desc.FindExcludeThirdPartyProtoFiles(pluginProtoDir)
		if err != nil {
			return err
		}

		for _, fp := range pluginProtoFiles {
			parser := rpcparser.NewDefaultProtoParser()
			parse, err := parser.Parse(fp, true)
			if err != nil {
				continue
			}

			for _, service := range parse.Service {
				pluginServices = append(pluginServices, service.Name)
			}

			var importPaths []string
			importPaths = append(importPaths, pluginProtoDir)
			importPaths = append(importPaths, pluginThirdPartyProtoDir)

			var protoParser protoparse.Parser
			protoParser.InferImportPaths = false

			protoParser.ImportPaths = []string{pluginProtoDir, pluginThirdPartyProtoDir}

			for _, v := range config.C.Gen.Zrpcclient.ProtoInclude {
				protoParser.ImportPaths = append(protoParser.ImportPaths, v)
			}
			protoParser.IncludeSourceCodeInfo = true

			protocCmd := fmt.Sprintf("protoc %s -I%s -I%s --go_out=%s --go-grpc_out=%s",
				fp,
				pluginProtoDir,
				pluginThirdPartyProtoDir,
				config.C.Gen.Zrpcclient.Output,
				config.C.Gen.Zrpcclient.Output,
			)

			for _, exp := range excludeThirdPartyProtoFiles {
				var expRel string
				var expErr error

				expRel, expErr = filepath.Rel(pluginProtoDir, exp)
				if expErr != nil {
					continue
				}

				var parserImportPaths []string
				parserImportPaths = []string{pluginProtoDir, pluginThirdPartyProtoDir}
				protoParser.ImportPaths = parserImportPaths

				fds, err := protoParser.ParseFiles(expRel)
				if err != nil {
					continue
				}

				if len(fds) == 0 {
					continue
				}

				expGoPackage := fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

				protocCmd += fmt.Sprintf(" --go_opt=module=%s --go_opt=M%s=%s --go-grpc_opt=module=%s --go-grpc_opt=M%s=%s", config.C.Gen.Zrpcclient.GoModule, expRel, func() string {
					if strings.HasPrefix(expGoPackage, config.C.Gen.Zrpcclient.GoModule) {
						return expGoPackage
					}
					return filepath.ToSlash(filepath.Join(config.C.Gen.Zrpcclient.GoModule, "plugins", p.Name, "model", expGoPackage))
				}(), config.C.Gen.Zrpcclient.GoModule, expRel, func() string {
					if strings.HasPrefix(expGoPackage, config.C.Gen.Zrpcclient.GoModule) {
						return expGoPackage
					}
					return filepath.ToSlash(filepath.Join(config.C.Gen.Zrpcclient.GoModule, "plugins", p.Name, "model", expGoPackage))
				}())
			}

			if len(config.C.Gen.Zrpcclient.ProtoInclude) > 0 {
				protocCmd += fmt.Sprintf(" -I%s ", strings.Join(config.C.Gen.Zrpcclient.ProtoInclude, " -I"))
			}

			logx.Debugf(protocCmd)
			resp, err := execx.Run(protocCmd, wd)
			if err != nil {
				return errors.Errorf("err: [%v], resp: [%s]", err, resp)
			}
		}

		if len(pluginServices) == 0 {
			continue
		}

		g := generator.NewGenerator(config.C.Style, false)

		for _, fp := range pluginProtoFiles {
			parser := rpcparser.NewDefaultProtoParser()
			parse, err := parser.Parse(fp, true)
			if err != nil {
				continue
			}

			pluginDirContext := DirContext{
				ImportBase:      filepath.Join(config.C.Gen.Zrpcclient.GoModule, "plugins", p.Name),
				PbPackage:       parse.PbPackage,
				OptionGoPackage: parse.GoPackage,
				Output:          filepath.Join(config.C.Gen.Zrpcclient.Output, "plugins", p.Name),
			}

			for _, service := range parse.Service {
				_ = os.MkdirAll(filepath.Join(pluginDirContext.GetCall().Filename, strings.ToLower(service.Name)), 0o755)
			}

			err = g.GenCall(pluginDirContext, parse, &conf.Config{
				NamingFormat: config.C.Style,
			}, &generator.ZRpcContext{
				Multiple:    true,
				IsGenClient: true,
			})
			if err != nil {
				return err
			}
		}

		pluginNames = append(pluginNames, p.Name)

		pluginTemplate, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugin.go.tpl")), map[string]any{
			"Module":     config.C.Gen.Zrpcclient.GoModule,
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

		pluginDir := filepath.Join(config.C.Gen.Zrpcclient.Output, "plugins")
		err = os.MkdirAll(pluginDir, 0o755)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(pluginDir, p.Name+".go"), formated, 0o644)
		if err != nil {
			return err
		}
	}

	pluginsTemplate, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugins.go.tpl")), map[string]any{
		"Module":      config.C.Gen.Zrpcclient.GoModule,
		"PluginNames": pluginNames,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "plugins.go.tpl"))))
	if err != nil {
		return err
	}

	formated, err := gosimports.Process("", pluginsTemplate, nil)
	if err != nil {
		return errors.Errorf("format plugins go file meet error: %v", err)
	}

	pluginDir := filepath.Join(config.C.Gen.Zrpcclient.Output, "plugins")
	err = os.MkdirAll(pluginDir, 0o755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(pluginDir, "plugins.go"), formated, 0o644)
	if err != nil {
		return err
	}

	for _, p := range plugins {
		if !pathx.FileExists(filepath.Join(p.Path, "desc", "proto")) {
			continue
		}

		err := genPluginNoRpcServiceExcludeThirdPartyProto(p, config.C.Gen.Zrpcclient.GoModule, config.C.Gen.Zrpcclient.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

func genPluginNoRpcServiceExcludeThirdPartyProto(plugin plugin.Plugin, goModule, output string) error {
	pluginProtoDir := filepath.Join(plugin.Path, "desc", "proto")
	pluginThirdPartyProtoDir := filepath.Join(plugin.Path, "desc", "proto", "third_party")

	excludeThirdPartyProtoFiles, err := desc.FindNoRpcServiceExcludeThirdPartyProtoFiles(pluginProtoDir)
	if err != nil || len(excludeThirdPartyProtoFiles) == 0 {
		return nil
	}

	var protoParser protoparse.Parser
	protoParser.InferImportPaths = false
	protoParser.ImportPaths = []string{pluginProtoDir, pluginThirdPartyProtoDir}
	protoParser.IncludeSourceCodeInfo = true

	pbDir := filepath.Join(output, "plugins", plugin.Name, "model")
	err = os.MkdirAll(pbDir, 0o755)
	if err != nil {
		return err
	}

	for _, v := range excludeThirdPartyProtoFiles {
		rel, err := filepath.Rel(pluginProtoDir, v)
		if err != nil {
			return err
		}

		fds, err := protoParser.ParseFiles(rel)
		if err != nil {
			continue
		}

		if len(fds) == 0 {
			continue
		}

		goPackage := fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

		command := fmt.Sprintf("protoc %s -I%s -I%s --go_out=%s --go_opt=module=%s --go_opt=M%s=%s --go-grpc_out=%s --go-grpc_opt=module=%s --go-grpc_opt=M%s=%s",
			v,
			pluginProtoDir,
			pluginThirdPartyProtoDir,
			output,
			goModule,
			rel,
			func() string {
				if strings.HasPrefix(goPackage, goModule) {
					return goPackage
				}
				return filepath.ToSlash(filepath.Join(goModule, "plugins", plugin.Name, "model", goPackage))
			}(),
			output,
			goModule,
			rel,
			func() string {
				if strings.HasPrefix(goPackage, goModule) {
					return goPackage
				}
				return filepath.ToSlash(filepath.Join(goModule, "plugins", plugin.Name, "model", goPackage))
			}(),
		)

		if len(config.C.Gen.Zrpcclient.ProtoInclude) > 0 {
			command += fmt.Sprintf(" -I%s ", strings.Join(config.C.Gen.Zrpcclient.ProtoInclude, " -I"))
		}

		logx.Debug(command)

		_, err = execx.Run(command, config.C.Wd())
		if err != nil {
			return err
		}
	}
	return nil
}

func genNoRpcServiceExcludeThirdPartyProto(genModule bool, module string) error {
	excludeThirdPartyProtoFiles, err := desc.FindNoRpcServiceExcludeThirdPartyProtoFiles(config.C.ProtoDir())
	if err != nil {
		return err
	}

	var protoParser protoparse.Parser
	protoParser.InferImportPaths = false

	protoDir := filepath.Join("desc", "proto")
	thirdPartyProtoDir := filepath.Join("desc", "proto", "third_party")
	protoParser.ImportPaths = []string{protoDir, thirdPartyProtoDir}
	protoParser.IncludeSourceCodeInfo = true

	pbDir := filepath.Join(config.C.Gen.Zrpcclient.Output, "model")
	err = os.MkdirAll(pbDir, 0o755)
	if err != nil {
		return err
	}

	for _, v := range excludeThirdPartyProtoFiles {
		rel, err := filepath.Rel(config.C.ProtoDir(), v)
		if err != nil {
			return err
		}

		fds, err := protoParser.ParseFiles(rel)
		if err != nil {
			return err
		}

		if len(fds) == 0 {
			continue
		}

		goPackage := fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

		getMod, err := mod.GetGoMod(config.C.Wd())
		if err != nil {
			return err
		}

		if !genModule {
			if config.C.Gen.Zrpcclient.Output != "." {
				module = getMod.Path
			}
		}

		command := fmt.Sprintf("protoc %s -I%s -I%s --go_out=%s --go_opt=module=%s --go_opt=M%s=%s --go-grpc_out=%s --go-grpc_opt=module=%s --go-grpc_opt=M%s=%s",
			v,
			config.C.ProtoDir(),
			filepath.Join(config.C.ProtoDir(), "third_party"),
			func() string {
				if !genModule {
					return "."
				}
				return filepath.Join(config.C.Gen.Zrpcclient.Output)
			}(),
			module,
			rel,
			func() string {
				if strings.HasPrefix(goPackage, module) {
					return goPackage
				}
				if genModule {
					return filepath.ToSlash(filepath.Join(module, "model", goPackage))
				}
				return filepath.ToSlash(filepath.Join(module, config.C.Gen.Zrpcclient.Output, "model", goPackage))
			}(),
			func() string {
				if !genModule {
					return "."
				}
				return filepath.Join(config.C.Gen.Zrpcclient.Output)
			}(),
			module,
			rel,
			func() string {
				if strings.HasPrefix(goPackage, module) {
					return goPackage
				}

				if genModule {
					return filepath.ToSlash(filepath.Join(module, "model", goPackage))
				}
				return filepath.ToSlash(filepath.Join(module, config.C.Gen.Zrpcclient.Output, "model", goPackage))
			}(),
		)

		logx.Debug(command)

		_, err = execx.Run(command, config.C.Wd())
		if err != nil {
			return err
		}
	}
	return nil
}
