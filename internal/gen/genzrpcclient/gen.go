package genzrpcclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/new"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type DirContext struct {
	ImportBase      string
	PbPackage       string
	OptionGoPackage string
	Scope           string
	Output          string
}

func (d DirContext) GetCall() generator.Dir {
	return generator.Dir{
		Filename: filepath.Join(d.Output, "typed", d.Scope),
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
		Package: filepath.ToSlash(fmt.Sprintf("%s/model/%s/%s", d.ImportBase, d.Scope, strings.TrimPrefix(d.OptionGoPackage, "./"))),
	}
}

func (d DirContext) GetProtoGo() generator.Dir {
	return generator.Dir{
		Filename: d.OptionGoPackage,
		Package:  filepath.ToSlash(fmt.Sprintf("%s/model/%s/%s", d.ImportBase, d.Scope, strings.TrimPrefix(d.OptionGoPackage, "./"))),
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

func Generate(c config.Config, genModule bool) error {
	g := generator.NewGenerator(c.Gen.Style, false)

	baseProtoDir := filepath.Join("desc", "proto")

	fps, err := desc.GetProtoFilepath(baseProtoDir)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	var services []string

	for _, fp := range fps {
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(fp, true)
		if err != nil {
			return err
		}
		dirContext := DirContext{
			ImportBase:      filepath.Join(c.Gen.Zrpcclient.GoModule),
			PbPackage:       parse.PbPackage,
			OptionGoPackage: parse.GoPackage,
			Scope:           c.Gen.Zrpcclient.Scope,
			Output:          c.Gen.Zrpcclient.Output,
		}
		for _, service := range parse.Service {
			services = append(services, service.Name)
			_ = os.MkdirAll(filepath.Join(dirContext.GetCall().Filename, strings.ToLower(service.Name)), 0o755)
		}

		// gen pb model
		err = os.MkdirAll(filepath.Join(c.Gen.Zrpcclient.Output, "model", c.Gen.Zrpcclient.Scope), 0o755)
		if err != nil {
			return err
		}
		resp, err := execx.Run(fmt.Sprintf("protoc -I%s -I%s --go_out=%s --go-grpc_out=%s %s", baseProtoDir, filepath.Join(baseProtoDir, "third_party"), filepath.Join(c.Gen.Zrpcclient.Output, "model", c.Gen.Zrpcclient.Scope), filepath.Join(c.Gen.Zrpcclient.Output, "model", c.Gen.Zrpcclient.Scope), fp), wd)
		if err != nil {
			return errors.Errorf("err: [%v], resp: [%s]", err, resp)
		}

		err = g.GenCall(dirContext, parse, &conf.Config{
			NamingFormat: c.Gen.Style,
		}, &generator.ZRpcContext{
			Multiple:    true,
			IsGenClient: true,
		})
		if err != nil {
			return err
		}
	}

	// gen clientset and options
	template, err := templatex.ParseTemplate(map[string]any{
		"Module":  c.Gen.Zrpcclient.GoModule,
		"Package": c.Gen.Zrpcclient.GoPackage,
		"Scopes":  []string{c.Gen.Zrpcclient.Scope},
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "clientset.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(c.Gen.Zrpcclient.Output, "clientset.go"), template, 0o644)
	if err != nil {
		return err
	}

	template, err = templatex.ParseTemplate(map[string]any{
		"Module":  c.Gen.Zrpcclient.GoModule,
		"Package": c.Gen.Zrpcclient.GoPackage,
		"Scopes":  []string{c.Gen.Zrpcclient.Scope},
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "options.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(c.Gen.Zrpcclient.Output, "options.go"), template, 0o644)
	if err != nil {
		return err
	}

	// generate scope client
	template, err = templatex.ParseTemplate(map[string]any{
		"Module":   c.Gen.Zrpcclient.GoModule,
		"Scope":    c.Gen.Zrpcclient.Scope,
		"Services": services,
	}, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "typed", "scope_client.go.tpl"))))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(c.Gen.Zrpcclient.Output, "typed", c.Gen.Zrpcclient.Scope, fmt.Sprintf("%s_client.go", c.Gen.Zrpcclient.Scope)), template, 0o644)
	if err != nil {
		return err
	}

	// if set --module flag
	if genModule {
		data, err := new.NewTemplateData()
		if err != nil {
			return err
		}
		data["Module"] = c.Gen.Zrpcclient.GoModule
		if c.Gen.Zrpcclient.GoVersion != "" {
			data["GoVersion"] = c.Gen.Zrpcclient.GoVersion
		}
		template, err = templatex.ParseTemplate(data, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "zrpcclient-go", "go.mod.tpl"))))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(c.Gen.Zrpcclient.Output, "go.mod"), template, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}
