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
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/new"
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

	// gen clientset and options
	template, err := templatex.ParseTemplate(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "clientset.go.tpl")), map[string]any{
		"Module":   config.C.Gen.Zrpcclient.GoModule,
		"Package":  config.C.Gen.Zrpcclient.GoPackage,
		"Services": services,
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

	return nil
}
