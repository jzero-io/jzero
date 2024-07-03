package genrpcclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/internal/gen"
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
	Output string
)

type DirContext struct {
	ImportBase      string
	PbPackage       string
	OptionGoPackage string
}

func (d DirContext) GetCall() generator.Dir {
	return generator.Dir{
		Filename: Output,
		GetChildPackage: func(childPath string) (string, error) {
			return childPath, nil
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
		Package: filepath.ToSlash(fmt.Sprintf("%s/%s", d.ImportBase, strings.TrimPrefix(d.OptionGoPackage, "./"))),
	}
}

func (d DirContext) GetProtoGo() generator.Dir {
	return generator.Dir{
		Filename: d.OptionGoPackage,
		Package:  filepath.ToSlash(fmt.Sprintf("%s/%s", d.ImportBase, strings.TrimPrefix(d.OptionGoPackage, "./"))),
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

	for _, fp := range fps {
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(fp, true)
		if err != nil {
			return err
		}
		dirContext := DirContext{
			ImportBase:      filepath.Join(module.Path, Output),
			PbPackage:       parse.PbPackage,
			OptionGoPackage: parse.GoPackage,
		}
		for _, service := range parse.Service {
			_ = os.MkdirAll(filepath.Join(dirContext.GetCall().Filename, service.Name), 0o755)
		}

		// gen pb model
		resp, err := execx.Run(fmt.Sprintf("protoc -I%s --go_out=%s --go-grpc_out=%s %s", baseProtoDir, dirContext.GetCall().Filename, dirContext.GetCall().Filename, fp), wd)
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

	return nil
}
