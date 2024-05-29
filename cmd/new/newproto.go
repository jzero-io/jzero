package new

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

type JzeroProto struct {
	TemplateData map[string]interface{}
}

func (jp *JzeroProto) New() error {
	protoDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app", "desc", "proto"))
	for _, file := range protoDir {
		if file.IsDir() {
			continue
		}
		protoFileBytes, err := templatex.ParseTemplate(jp.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "proto", file.Name())))
		cobra.CheckErr(err)
		protoFileName := file.Name()
		err = checkWrite(filepath.Join(Output, "app", "desc", "proto", protoFileName), protoFileBytes)
		cobra.CheckErr(err)

		if len(protoFileBytes) > 0 {
			if !pathx.FileExists(filepath.Join(Output, "app", "desc", "proto", "google", "protobuf")) {
				err = embeded.WriteTemplateDir(filepath.Join("jzero", "app", "desc", "proto", "google", "protobuf"), filepath.Join(Output, "app", "desc", "proto", "google", "protobuf"))
				cobra.CheckErr(err)
			}
		}
	}
	if isNeedNewGoogleApiProto(filepath.Join(Output, "app", "desc", "proto")) {
		err := embeded.WriteTemplateDir(filepath.Join("jzero", "app", "desc", "proto", "google", "api"), filepath.Join(Output, "app", "desc", "proto", "google", "api"))
		if err != nil {
			return err
		}
	}
	return nil
}

func isNeedNewGoogleApiProto(protoDirName string) bool {
	protoDir, err := os.ReadDir(protoDirName)
	if err != nil {
		return false
	}

	for _, v := range protoDir {
		if v.IsDir() {
			continue
		}
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(filepath.Join(protoDirName, v.Name()))
		if err != nil {
			return false
		}
		for _, ps := range parse.Service {
			for _, rpc := range ps.RPC {
				for _, option := range rpc.Options {
					if option.Name == "(google.api.http)" {
						return true
					}
				}
			}
		}
	}
	return false
}
