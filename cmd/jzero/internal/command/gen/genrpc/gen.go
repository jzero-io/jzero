package genrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gitstatus"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
)

type JzeroRpc struct {
	Module string

	ProtoFiles          []string
	GenCodeProtoFiles   []string
	ProtoSpecMap        map[string]rpcparser.Proto
	GenCodeProtoSpecMap map[string]rpcparser.Proto
}

type (
	ImportLines []string

	RegisterLines []string
)

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

func (jr *JzeroRpc) Gen() error {
	var (
		serverImports   ImportLines
		pbImports       ImportLines
		registerServers RegisterLines
	)

	// 获取全量 proto 文件
	protoFiles, err := jzerodesc.GetProtoFilepath(config.C.ProtoDir())
	if err != nil {
		return err
	}
	jr.ProtoFiles = protoFiles
	if len(jr.ProtoFiles) == 0 {
		return nil
	}

	jr.ProtoSpecMap = make(map[string]rpcparser.Proto, len(protoFiles))
	for _, v := range protoFiles {
		// parse proto
		protoParser := rpcparser.NewDefaultProtoParser()
		var parse rpcparser.Proto
		parse, err = protoParser.Parse(v, true)
		if err != nil {
			return err
		}
		jr.ProtoSpecMap[v] = parse
	}

	// 获取需要生成代码的proto 文件
	var genCodeProtoFiles []string
	jr.GenCodeProtoSpecMap = make(map[string]rpcparser.Proto, len(protoFiles))

	switch {
	case config.C.Gen.GitChange && filex.DirExists(filepath.Join(config.C.Wd(), ".git")) && len(config.C.Gen.Desc) == 0:
		m, _, err := gitstatus.ChangedFiles(config.C.ProtoDir(), ".proto")
		if err == nil {
			genCodeProtoFiles = append(genCodeProtoFiles, m...)
			for _, file := range m {
				jr.GenCodeProtoSpecMap[file] = jr.ProtoSpecMap[file]
			}
		}
	case len(config.C.Gen.Desc) > 0:
		for _, v := range config.C.Gen.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					genCodeProtoFiles = append(genCodeProtoFiles, filepath.Join(strings.Split(filepath.ToSlash(v), "/")...))
					jr.GenCodeProtoSpecMap[filepath.Clean(v)] = jr.ProtoSpecMap[filepath.Clean(v)]
				}
			} else {
				specifiedProtoFiles, err := jzerodesc.GetProtoFilepath(v)
				if err != nil {
					return err
				}
				genCodeProtoFiles = append(genCodeProtoFiles, specifiedProtoFiles...)
				for _, saf := range specifiedProtoFiles {
					jr.GenCodeProtoSpecMap[filepath.Clean(saf)] = jr.ProtoSpecMap[filepath.Clean(saf)]
				}
			}
		}
	default:
		// 否则生成代码的 proto 文件为全量 proto 文件
		genCodeProtoFiles = jr.ProtoFiles
		jr.GenCodeProtoSpecMap = jr.ProtoSpecMap
	}
	jr.GenCodeProtoFiles = genCodeProtoFiles

	// ignore proto desc
	for _, v := range config.C.Gen.DescIgnore {
		if !osx.IsDir(v) {
			if filepath.Ext(v) == ".proto" {
				// delete item in genCodeApiFiles by filename
				jr.GenCodeProtoFiles = lo.Reject(jr.GenCodeProtoFiles, func(item string, _ int) bool {
					return item == v
				})
				// delete map key
				delete(jr.GenCodeProtoSpecMap, v)
			}
		} else {
			specifiedApiFiles, err := jzerodesc.GetProtoFilepath(v)
			if err != nil {
				return err
			}
			for _, saf := range specifiedApiFiles {
				jr.GenCodeProtoFiles = lo.Reject(jr.GenCodeProtoFiles, func(item string, _ int) bool {
					return item == saf
				})
				delete(jr.GenCodeProtoSpecMap, saf)
			}
		}
	}

	if len(jr.GenCodeProtoFiles) == 0 {
		return nil
	}

	fmt.Printf("%s to generate proto code. \n", color.WithColor("Start", color.FgGreen))
	if err = jr.patchSvc(); err != nil {
		return err
	}

	for _, v := range jr.ProtoFiles {
		allLogicFiles, err := jr.GetAllLogicFiles(v, jr.ProtoSpecMap[v])
		if err != nil {
			return err
		}

		allServerFiles, err := jr.GetAllServerFiles(v, jr.ProtoSpecMap[v])
		if err != nil {
			return err
		}

		if config.C.Gen.RpcStylePatch {
			if lo.Contains(genCodeProtoFiles, v) {
				for _, s := range jr.ProtoSpecMap[v].Service {
					// rename logic dir and server dir
					dirName, _ := format.FileNamingFormat("gozero", s.Name)
					fixDirName, _ := format.FileNamingFormat(config.C.Gen.Style, s.Name)

					_ = os.Rename(filepath.Join("internal", "logic", strings.ToLower(fixDirName)), filepath.Join("internal", "logic", dirName))
					_ = os.Rename(filepath.Join("internal", "server", strings.ToLower(fixDirName)), filepath.Join("internal", "server", dirName))
				}
			}
		}

		if lo.Contains(genCodeProtoFiles, v) {
			fmt.Printf("%s proto file %s\n", color.WithColor("Using", color.FgGreen), v)
			zrpcOut := "."
			command := fmt.Sprintf("goctl rpc protoc %s -I%s -I%s --go_out=%s --go-grpc_out=%s --zrpc_out=%s --client=%t --home %s -m --style %s ",
				v,
				config.C.ProtoDir(),
				filepath.Join(config.C.ProtoDir(), "third_party"),
				filepath.Join("internal"),
				filepath.Join("internal"),
				zrpcOut,
				config.C.Gen.RpcClient,
				filepath.Join(embeded.Home, "go-zero"),
				config.C.Gen.Style)

			logx.Debug(command)

			_, err = execx.Run(command, config.C.Wd())
			if err != nil {
				return err
			}
		}

		if config.C.Gen.RpcStylePatch {
			if lo.Contains(genCodeProtoFiles, v) {
				for _, s := range jr.ProtoSpecMap[v].Service {
					// rename logic dir and server dir
					dirName, _ := format.FileNamingFormat("gozero", s.Name)
					fixDirName, _ := format.FileNamingFormat(config.C.Gen.Style, s.Name)

					_ = os.Rename(filepath.Join("internal", "logic", strings.ToLower(dirName)), filepath.Join("internal", "logic", fixDirName))
					_ = os.Rename(filepath.Join("internal", "server", strings.ToLower(dirName)), filepath.Join("internal", "server", fixDirName))
				}
			}
		}

		for _, file := range allServerFiles {
			if filepath.Clean(file.DescFilepath) == filepath.Clean(v) {
				if _, ok := jr.GenCodeProtoSpecMap[file.DescFilepath]; ok {
					if err := jr.removeServerSuffix(file.Path); err != nil {
						console.Warning("[warning]: remove server suffix %s meet error %v", file.Path, err)
						continue
					}
				}
			}
		}
		for _, file := range allLogicFiles {
			if _, ok := jr.GenCodeProtoSpecMap[file.DescFilepath]; ok {
				if err := jr.removeLogicSuffix(file.Path); err != nil {
					console.Warning("[warning]: remove logic suffix %s meet error %v", file.Path, err)
					continue
				}
			}
		}

		if config.C.Gen.RpcStylePatch {
			if lo.Contains(genCodeProtoFiles, v) {
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
		}

		if config.C.Gen.ChangeLogicTypes {
			if lo.Contains(genCodeProtoFiles, v) {
				for _, file := range allLogicFiles {
					if err = jr.changeLogicTypes(file); err != nil {
						console.Warning("[warning]: change logic types %s meet error %v", file.Path, err)
						continue
					}
				}
			}
		}

		// # gen proto descriptor
		if lo.Contains(genCodeProtoFiles, v) {
			if isNeedGenProtoDescriptor(jr.ProtoSpecMap[v]) {
				if !pathx.FileExists(getProtoDescriptorPath(v)) {
					_ = os.MkdirAll(filepath.Dir(getProtoDescriptorPath(v)), 0o755)
				}
				protocCommand := fmt.Sprintf("protoc --include_imports -I%s -I%s --descriptor_set_out=%s %s",
					config.C.ProtoDir(),
					filepath.Join(config.C.ProtoDir(), "third_party"),
					getProtoDescriptorPath(v),
					v,
				)
				_, err = execx.Run(protocCommand, config.C.Wd())
				if err != nil {
					return err
				}
			}
		}

		for _, s := range jr.ProtoSpecMap[v].Service {
			if config.C.Gen.RpcStylePatch {
				serverDir, _ := format.FileNamingFormat(config.C.Gen.Style, s.Name)
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, strings.ToLower(s.Name), jr.Module, strings.ToLower(serverDir)))
			} else {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, strings.ToLower(s.Name), jr.Module, strings.ToLower(s.Name)))
			}

			registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%s(ctx))", filepath.Base(jr.ProtoSpecMap[v].GoPackage), stringx.FirstUpper(s.Name), strings.ToLower(s.Name), stringx.FirstUpper(stringx.ToCamel(s.Name))))
		}
		pbImports = append(pbImports, fmt.Sprintf(`"%s/internal/%s"`, jr.Module, strings.TrimPrefix(jr.ProtoSpecMap[v].GoPackage, "./")))
	}
	if len(jr.GenCodeProtoFiles) > 0 {
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	if pathx.FileExists(config.C.ProtoDir()) {
		if err = jr.genServer(serverImports, pbImports, registerServers); err != nil {
			return err
		}
		if err = jr.genApiMiddlewares(); err != nil {
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
