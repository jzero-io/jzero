package genrpc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
)

func (jr *JzeroRpc) genNoRpcServiceExcludeThirdPartyProto(protoDirPath string) error {
	excludeThirdPartyProtoFiles, err := jzerodesc.FindNoRpcServiceExcludeThirdPartyProtoFiles(protoDirPath)
	if err != nil {
		return err
	}

	// 获取 proto 的 go_package
	var protoParser protoparse.Parser
	protoParser.InferImportPaths = false

	protoDir := filepath.Join("desc", "proto")
	thirdPartyProtoDir := filepath.Join("desc", "proto", "third_party")
	protoParser.ImportPaths = []string{protoDir, thirdPartyProtoDir}
	protoParser.IncludeSourceCodeInfo = true

	for _, v := range excludeThirdPartyProtoFiles {
		rel, err := filepath.Rel(protoDir, v)
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

		command := fmt.Sprintf("protoc %s -I%s -I%s --go_out=%s --go_opt=module=%s --go_opt=M%s=%s --go-grpc_out=%s --go-grpc_opt=module=%s",
			v,
			config.C.ProtoDir(),
			filepath.Join(config.C.ProtoDir(), "third_party"),
			filepath.Join("."),
			jr.Module,
			rel,
			func() string {
				if strings.HasPrefix(goPackage, jr.Module) {
					return goPackage
				}
				return filepath.ToSlash(filepath.Join(jr.Module, "internal", goPackage))
			}(),
			filepath.Join("."),
			jr.Module)

		logx.Debug(command)

		_, err = execx.Run(command, config.C.Wd())
		if err != nil {
			return err
		}
	}
	return nil
}
