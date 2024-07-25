package genzrpcclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"

	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/internal/new"
	"github.com/jzero-io/jzero/pkg/mod"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

var (
	Style  string
	Scope  string
	Output string
	Module string
)

type DirContext struct {
	ImportBase      string
	PbPackage       string
	OptionGoPackage string
}

func (d DirContext) GetCall() generator.Dir {
	return generator.Dir{
		Filename: filepath.Join(Output, "typed", Scope),
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
		Package: filepath.ToSlash(fmt.Sprintf("%s/model/%s/%s", d.ImportBase, Scope, strings.TrimPrefix(d.OptionGoPackage, "./"))),
	}
}

func (d DirContext) GetProtoGo() generator.Dir {
	return generator.Dir{
		Filename: d.OptionGoPackage,
		Package:  filepath.ToSlash(fmt.Sprintf("%s/model/%s/%s", d.ImportBase, Scope, strings.TrimPrefix(d.OptionGoPackage, "./"))),
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

func Generate(_ *cobra.Command, _ []string) error {
	g := generator.NewGenerator(Style, false)

	baseProtoDir := filepath.Join("desc", "proto")

	fps, err := gen.GetProtoFilepath(baseProtoDir)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	module, err := mod.GetGoMod(wd)
	if err != nil {
		return err
	}
	if Module != "" {
		module.Path = Module
	} else {
		module.Path = filepath.ToSlash(filepath.Join(module.Path, Output))
	}

	var services []string

	for _, fp := range fps {
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(fp, true)
		if err != nil {
			return err
		}
		dirContext := DirContext{
			ImportBase:      filepath.Join(module.Path),
			PbPackage:       parse.PbPackage,
			OptionGoPackage: parse.GoPackage,
		}
		for _, service := range parse.Service {
			services = append(services, service.Name)
			_ = os.MkdirAll(filepath.Join(dirContext.GetCall().Filename, strings.ToLower(service.Name)), 0o755)
		}

		// gen pb model
		err = os.MkdirAll(filepath.Join(Output, "model", Scope), 0o755)
		if err != nil {
			return err
		}
		resp, err := execx.Run(fmt.Sprintf("protoc -I%s --go_out=%s --go-grpc_out=%s %s", baseProtoDir, filepath.Join(Output, "model", Scope), filepath.Join(Output, "model", Scope), fp), wd)
		if err != nil {
			return errors.Errorf("err: [%v], resp: [%s]", err, resp)
		}

		err = g.GenCall(dirContext, parse, &conf.Config{
			NamingFormat: Style,
		}, &generator.ZRpcContext{
			Multiple:    true,
			IsGenClient: true,
		})
		if err != nil {
			return err
		}
	}

	// gen clientset and options
	template, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": module.Path,
		"APP":    Scope,
		"Scopes": []string{Scope},
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "clientset.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(Output, "clientset.go"), template, 0o644)
	if err != nil {
		return err
	}

	template, err = templatex.ParseTemplate(map[string]interface{}{
		"Module": module.Path,
		"APP":    Scope,
		"Scopes": []string{Scope},
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "options.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(Output, "options.go"), template, 0o644)
	if err != nil {
		return err
	}

	// generate scope client
	template, err = templatex.ParseTemplate(map[string]interface{}{
		"Module":   module.Path,
		"Scope":    Scope,
		"Services": services,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "typed", "scope_client.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(Output, "typed", Scope, fmt.Sprintf("%s_client.go", Scope)), template, 0o644)
	if err != nil {
		return err
	}

	// if set --module flag
	if Module != "" {
		data, err := new.NewTemplateData(nil)
		if err != nil {
			return err
		}
		data["Module"] = Module
		template, err = templatex.ParseTemplate(data, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl"))))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(Output, "go.mod"), template, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}
