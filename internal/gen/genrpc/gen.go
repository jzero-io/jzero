package genrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen/genapi"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/stringx"
)

type JzeroRpc struct {
	Wd     string
	Module string
	config.GenConfig
}

func (jr *JzeroRpc) Gen() error {
	protoDirPath := filepath.Join("desc", "proto")
	protoFilenames, err := jzerodesc.GetProtoFilepath(protoDirPath)
	if err != nil {
		return err
	}

	var serverImports jzerodesc.ImportLines
	var pbImports jzerodesc.ImportLines
	var registerServers jzerodesc.RegisterLines
	var allServerFiles []ServerFile
	var allLogicFiles []genapi.LogicFile

	if len(protoFilenames) > 0 {
		fmt.Printf("%s to generate proto code. \n", color.WithColor("Start", color.FgGreen))
	}

	for _, v := range protoFilenames {
		// parse proto
		protoParser := rpcparser.NewDefaultProtoParser()
		var parse rpcparser.Proto
		parse, err = protoParser.Parse(v, true)
		if err != nil {
			return err
		}

		allLogicFiles, err = jr.GetAllLogicFiles(parse)
		if err != nil {
			return err
		}

		allServerFiles, err = jr.GetAllServerFiles(parse)
		if err != nil {
			return err
		}

		if jr.RpcStylePatch {
			for _, s := range parse.Service {
				// rename logic dir and server dir
				dirName, _ := format.FileNamingFormat("gozero", s.Name)
				fixDirName, _ := format.FileNamingFormat(jr.Style, s.Name)

				_ = os.Rename(filepath.Join("internal", "logic", strings.ToLower(fixDirName)), filepath.Join("internal", "logic", dirName))
				_ = os.Rename(filepath.Join("internal", "server", strings.ToLower(fixDirName)), filepath.Join("internal", "server", dirName))
			}
		}

		fmt.Printf("%s proto file %s\n", color.WithColor("Using", color.FgGreen), v)
		zrpcOut := "."

		command := fmt.Sprintf("goctl rpc protoc %s -I%s -I%s --go_out=%s --go-grpc_out=%s --zrpc_out=%s --client=false --home %s -m --style %s ",
			v,
			protoDirPath,
			filepath.Join(protoDirPath, "third_party"),
			filepath.Join("internal"),
			filepath.Join("internal"),
			zrpcOut,
			filepath.Join(embeded.Home, "go-zero"),
			jr.Style)

		logx.Debug(command)

		_, err = execx.Run(command, jr.Wd)
		if err != nil {
			return err
		}

		if jr.RpcStylePatch {
			for _, s := range parse.Service {
				// rename logic dir and server dir
				dirName, _ := format.FileNamingFormat("gozero", s.Name)
				fixDirName, _ := format.FileNamingFormat(jr.Style, s.Name)

				_ = os.Rename(filepath.Join("internal", "logic", strings.ToLower(dirName)), filepath.Join("internal", "logic", fixDirName))
				_ = os.Rename(filepath.Join("internal", "server", strings.ToLower(dirName)), filepath.Join("internal", "server", fixDirName))
			}
		}

		command = fmt.Sprintf("protoc %s -I%s -I%s --validate_out=%s",
			v,
			protoDirPath,
			filepath.Join(protoDirPath, "third_party"),
			"lang=go:internal",
		)
		_, err = execx.Run(command, jr.Wd)
		if err != nil {
			return err
		}

		if jr.RemoveSuffix {
			for _, file := range allServerFiles {
				if err := jr.removeServerSuffix(file.Path); err != nil {
					console.Warning("[warning]: remove server suffix %s meet error %v", file.Path, err)
					continue
				}
			}
			for _, file := range allLogicFiles {
				if err := jr.removeLogicSuffix(file.Path); err != nil {
					console.Warning("[warning]: remove logic suffix %s meet error %v", file.Path, err)
					continue
				}
			}
		}

		if jr.RpcStylePatch {
			for _, file := range allServerFiles {
				err = jr.rpcStylePatchServer(file)
				if err != nil {
					return err
				}
			}
			for _, file := range allLogicFiles {
				err = jr.rpcStylePatchLogic(file)
				if err != nil {
					return err
				}
			}
		}

		if jr.ChangeLogicTypes {
			for _, file := range allLogicFiles {
				if err := jr.changeLogicTypes(file); err != nil {
					console.Warning("[warning]: change logic types %s meet error %v", file.Path, err)
					continue
				}
			}
		}

		// # gen proto descriptor
		if isNeedGenProtoDescriptor(parse) {
			if !pathx.FileExists(getProtoDescriptorPath(v)) {
				_ = os.MkdirAll(filepath.Dir(getProtoDescriptorPath(v)), 0o755)
			}
			protocCommand := fmt.Sprintf("protoc --include_imports -I%s -I%s --descriptor_set_out=%s %s",
				protoDirPath,
				filepath.Join(protoDirPath, "third_party"),
				getProtoDescriptorPath(v),
				v,
			)
			_, err = execx.Run(protocCommand, jr.Wd)
			if err != nil {
				return err
			}
		}

		for _, s := range parse.Service {
			if jr.RpcStylePatch {
				serverDir, _ := format.FileNamingFormat(jr.Style, s.Name)
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, strings.ToLower(s.Name), jr.Module, strings.ToLower(serverDir)))
			} else {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, strings.ToLower(s.Name), jr.Module, strings.ToLower(s.Name)))
			}

			if jr.RemoveSuffix {
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%s(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), strings.ToLower(s.Name), stringx.FirstUpper(stringx.ToCamel(s.Name))))
			} else {
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), strings.ToLower(s.Name), stringx.FirstUpper(stringx.ToCamel(s.Name))))
			}
		}
		pbImports = append(pbImports, fmt.Sprintf(`"%s/internal/%s"`, jr.Module, strings.TrimPrefix(parse.GoPackage, "./")))
	}
	if len(protoFilenames) > 0 {
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	if pathx.FileExists(protoDirPath) {
		if err = jr.genServer(serverImports, pbImports, registerServers); err != nil {
			return err
		}
		if err = jr.genApiMiddlewares(protoFilenames); err != nil {
			return err
		}
	}
	return nil
}

func getProtoDescriptorPath(protoPath string) string {
	rel, err := filepath.Rel(filepath.Join("desc", "proto"), protoPath)
	if err != nil {
		return ""
	}

	return filepath.Join("desc", "pb", strings.TrimSuffix(rel, ".proto")+".pb")
}

func isNeedGenProtoDescriptor(proto rpcparser.Proto) bool {
	for _, ps := range proto.Service {
		for _, rpc := range ps.RPC {
			for _, option := range rpc.Options {
				if option.Name == "(google.api.http)" {
					return true
				}
			}
		}
	}
	return false
}
