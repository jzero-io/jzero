package genzrpcclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"

	new2 "github.com/jzero-io/jzero/cmd/jzero/internal/command/new"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type DirContext struct {
	ImportBase      string
	PbPackage       string
	OptionGoPackage string
	Scope           string
	Output          string
	PbDir           string
	ClientDir       string
}

func (d DirContext) GetCall() generator.Dir {
	fileName := filepath.Join(d.Output, "typed", d.Scope)
	if d.ClientDir != "" {
		fileName = filepath.Join(d.Output, d.ClientDir)
	}
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
	packagePath := filepath.ToSlash(fmt.Sprintf("%s/model/%s/%s", d.ImportBase, d.Scope, strings.TrimPrefix(d.OptionGoPackage, "./")))
	if d.PbDir != "" {
		packagePath = filepath.ToSlash(fmt.Sprintf("%s/%s/%s", d.ImportBase, d.PbDir, strings.TrimPrefix(d.OptionGoPackage, "./")))
	}
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
			Scope:           config.C.Gen.Zrpcclient.Scope,
			Output:          config.C.Gen.Zrpcclient.Output,
			PbDir:           config.C.Gen.Zrpcclient.PbDir,
			ClientDir:       config.C.Gen.Zrpcclient.ClientDir,
		}
		for _, service := range parse.Service {
			services = append(services, service.Name)
			_ = os.MkdirAll(filepath.Join(dirContext.GetCall().Filename, strings.ToLower(service.Name)), 0o755)
		}
		pbDir := filepath.Join(config.C.Gen.Zrpcclient.Output, "model", config.C.Gen.Zrpcclient.Scope)
		if dirContext.PbDir != "" {
			pbDir = filepath.Join(config.C.Gen.Zrpcclient.Output, dirContext.PbDir)
		}
		// gen pb model
		err = os.MkdirAll(pbDir, 0o755)
		if err != nil {
			return err
		}
		resp, err := execx.Run(fmt.Sprintf("protoc -I%s -I%s --go_out=%s --go-grpc_out=%s %s", config.C.ProtoDir(), filepath.Join(config.C.ProtoDir(), "third_party"), pbDir, pbDir, fp), wd)
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

	var clientDir string
	filePath := filepath.Join(config.C.Gen.Zrpcclient.Output, "typed", config.C.Gen.Zrpcclient.Scope, fmt.Sprintf("%s_client.go", config.C.Gen.Zrpcclient.Scope))
	if config.C.Gen.Zrpcclient.ClientDir != "" {
		filePath = filepath.Join(config.C.Gen.Zrpcclient.Output, config.C.Gen.Zrpcclient.ClientDir, fmt.Sprintf("%s_client.go", config.C.Gen.Zrpcclient.Scope))
		clientDir = config.C.Gen.Zrpcclient.ClientDir
	} else {
		clientDir = "typed/" + config.C.Gen.Zrpcclient.Scope
	}

	// gen clientset and options
	template, err := templatex.ParseTemplate(map[string]any{
		"Module":    config.C.Gen.Zrpcclient.GoModule,
		"Package":   config.C.Gen.Zrpcclient.GoPackage,
		"Scopes":    []string{config.C.Gen.Zrpcclient.Scope},
		"ClientDir": clientDir,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "clientset.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(config.C.Gen.Zrpcclient.Output, "clientset.go"), template, 0o644)
	if err != nil {
		return err
	}

	template, err = templatex.ParseTemplate(map[string]any{
		"Module":    config.C.Gen.Zrpcclient.GoModule,
		"Package":   config.C.Gen.Zrpcclient.GoPackage,
		"Scopes":    []string{config.C.Gen.Zrpcclient.Scope},
		"ClientDir": clientDir,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "options.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(config.C.Gen.Zrpcclient.Output, "options.go"), template, 0o644)
	if err != nil {
		return err
	}

	// generate scope client
	scope := "typed/" + config.C.Gen.Zrpcclient.Scope
	if config.C.Gen.Zrpcclient.PbDir != "" {
		scope = config.C.Gen.Zrpcclient.ClientDir
	}
	template, err = templatex.ParseTemplate(map[string]any{
		"Module":   config.C.Gen.Zrpcclient.GoModule,
		"Scope":    scope,
		"Package":  config.C.Gen.Zrpcclient.Scope,
		"Services": services,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "typed", "scope_client.go.tpl"))))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, template, 0o644)
	if err != nil {
		return err
	}

	// if set --module flag
	if genModule {
		data, err := new2.NewTemplateData()
		if err != nil {
			return err
		}
		data["Module"] = config.C.Gen.Zrpcclient.GoModule
		if config.C.Gen.Zrpcclient.GoVersion != "" {
			data["GoVersion"] = config.C.Gen.Zrpcclient.GoVersion
		}
		template, err = templatex.ParseTemplate(data, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl"))))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(config.C.Gen.Zrpcclient.Output, "go.mod"), template, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}
