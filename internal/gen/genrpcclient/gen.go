package genrpcclient

import (
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/spf13/cobra"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"path/filepath"
)

var (
	Style string
)

type DirContext struct{}

func (d DirContext) GetCall() generator.Dir {
	return generator.Dir{
		Base:     "client",
		Filename: "",
		Package:  "",
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
	return generator.Dir{}
}

func (d DirContext) GetProtoGo() generator.Dir {
	return generator.Dir{}
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

func Generate(command *cobra.Command, args []string) error {
	g := generator.NewGenerator(Style, false)

	var dirContext DirContext

	fps, err := gen.GetProtoFilepath(filepath.Join("desc", "proto"))
	if err != nil {
		return err
	}
	for _, fp := range fps {
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(fp, true)
		if err != nil {
			return err
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
