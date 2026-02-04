package genrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	goctlconsole "github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gitstatus"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
)

type JzeroRpc struct {
	Module string
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

func (jr *JzeroRpc) Gen() (map[string]rpcparser.Proto, error) {
	var (
		serverImports   ImportLines
		pbImports       ImportLines
		registerServers RegisterLines
	)

	// 获取全量 proto 文件
	protoFiles, err := jzerodesc.FindRpcServiceProtoFiles(config.C.ProtoDir())
	if err != nil {
		return nil, err
	}

	if len(protoFiles) == 0 {
		return nil, nil
	}

	protoSpecMap := make(map[string]rpcparser.Proto, len(protoFiles))
	for _, v := range protoFiles {
		// parse proto
		protoParser := rpcparser.NewDefaultProtoParser()
		var parse rpcparser.Proto
		parse, err = protoParser.Parse(v, true)
		if err != nil {
			return nil, err
		}
		protoSpecMap[v] = parse
	}

	// 获取需要生成代码的proto 文件
	var genCodeProtoFiles []string
	genCodeProtoSpecMap := make(map[string]rpcparser.Proto, len(protoFiles))

	switch {
	case config.C.Gen.GitChange && gitstatus.IsGitRepo(filepath.Join(config.C.Wd())) && len(config.C.Gen.Desc) == 0:
		m, _, err := gitstatus.ChangedFiles(config.C.ProtoDir(), ".proto")
		if err == nil {
			genCodeProtoFiles = append(genCodeProtoFiles, m...)
			for _, file := range m {
				genCodeProtoSpecMap[file] = protoSpecMap[file]
			}
		}
	case len(config.C.Gen.Desc) > 0:
		for _, v := range config.C.Gen.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					genCodeProtoFiles = append(genCodeProtoFiles, filepath.Join(strings.Split(filepath.ToSlash(v), "/")...))
					genCodeProtoSpecMap[filepath.Clean(v)] = protoSpecMap[filepath.Clean(v)]
				}
			} else {
				specifiedProtoFiles, err := jzerodesc.FindRpcServiceProtoFiles(v)
				if err != nil {
					return nil, err
				}
				genCodeProtoFiles = append(genCodeProtoFiles, specifiedProtoFiles...)
				for _, saf := range specifiedProtoFiles {
					genCodeProtoSpecMap[filepath.Clean(saf)] = protoSpecMap[filepath.Clean(saf)]
				}
			}
		}
	default:
		// 否则生成代码的 proto 文件为全量 proto 文件
		genCodeProtoFiles = protoFiles
		genCodeProtoSpecMap = protoSpecMap
	}

	// ignore proto desc
	for _, v := range config.C.Gen.DescIgnore {
		if !osx.IsDir(v) {
			if filepath.Ext(v) == ".proto" {
				// delete item in genCodeApiFiles by filename
				genCodeProtoFiles = lo.Reject(genCodeProtoFiles, func(item string, _ int) bool {
					return item == filepath.Clean(v)
				})
				protoFiles = lo.Reject(protoFiles, func(item string, _ int) bool {
					return item == filepath.Clean(v)
				})
				// delete map key
				delete(genCodeProtoSpecMap, filepath.Clean(v))
				delete(protoSpecMap, filepath.Clean(v))
			}
		} else {
			specifiedProtoFiles, err := jzerodesc.FindRpcServiceProtoFiles(v)
			if err != nil {
				return nil, err
			}
			for _, saf := range specifiedProtoFiles {
				genCodeProtoFiles = lo.Reject(genCodeProtoFiles, func(item string, _ int) bool {
					return item == saf
				})
				protoFiles = lo.Reject(protoFiles, func(item string, _ int) bool {
					return item == saf
				})
				delete(genCodeProtoSpecMap, saf)
				delete(protoSpecMap, saf)
			}
		}
	}

	if len(genCodeProtoFiles) == 0 {
		return protoSpecMap, nil
	}

	if config.C.Quiet {
		fmt.Printf("%s to generate rpc code from proto files\n", console.Green("Start"))
	}

	// 处理模板
	var goctlHome string
	tempDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// // 先写入内置模板
	// err = embeded.WriteTemplateDir(filepath.Join("go-zero", "rpc"), filepath.Join(tempDir, "rpc"))
	// if err != nil {
	//	return nil, err
	// }

	// 如果用户自定义了模板，则复制覆盖
	customTemplatePath := filepath.Join(config.C.Home, "go-zero", "rpc")
	if pathx.FileExists(customTemplatePath) {
		err = filex.CopyDir(customTemplatePath, filepath.Join(tempDir, "rpc"))
		if err != nil {
			return nil, err
		}
	}

	goctlHome = tempDir
	logx.Debugf("goctl_home = %s", goctlHome)

	excludeThirdPartyProtoFiles, err := jzerodesc.FindExcludeThirdPartyProtoFiles(config.C.ProtoDir())
	if err != nil {
		return nil, err
	}
	logx.Debugf("excludeThirdPartyProtoFiles = %s", excludeThirdPartyProtoFiles)

	// 获取 proto 的 go_package
	var protoParser protoparse.Parser
	protoParser.InferImportPaths = false

	protoDir := filepath.Join("desc", "proto")
	thirdPartyProtoDir := filepath.Join("desc", "proto", "third_party")
	protoParser.ImportPaths = []string{protoDir, thirdPartyProtoDir}
	for _, v := range config.C.Gen.ProtoInclude {
		protoParser.ImportPaths = append(protoParser.ImportPaths, v)
	}
	protoParser.IncludeSourceCodeInfo = true

	for _, v := range protoFiles {
		allLogicFiles, err := jr.GetAllLogicFiles(v, protoSpecMap[v])
		if err != nil {
			return nil, err
		}

		allServerFiles, err := jr.GetAllServerFiles(v, protoSpecMap[v])
		if err != nil {
			return nil, err
		}

		if lo.Contains(genCodeProtoFiles, v) {
			if !config.C.Quiet {
				fmt.Printf("%s proto file %s\n", console.Green("Using"), v)
			}
			zrpcOut := "."

			rel, err := filepath.Rel(protoDir, v)
			if err != nil {
				return nil, err
			}

			fds, err := protoParser.ParseFiles(rel)
			if err != nil {
				return nil, err
			}

			if len(fds) == 0 {
				continue
			}

			goPackage := fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

			command := fmt.Sprintf("goctl rpc protoc %s -I%s -I%s --go_out=%s --go-grpc_out=%s --zrpc_out=%s --client=%t --home %s -m --style %s",
				v,
				config.C.ProtoDir(),
				filepath.Join(config.C.ProtoDir(), "third_party"),
				filepath.Join("."),
				filepath.Join("."),
				zrpcOut,
				config.C.Gen.RpcClient,
				goctlHome,
				config.C.Style)

			for _, exp := range excludeThirdPartyProtoFiles {
				rel, err = filepath.Rel(config.C.ProtoDir(), exp)
				if err != nil {
					return nil, err
				}

				fds, err = protoParser.ParseFiles(rel)
				if err != nil {
					return nil, err
				}

				if len(fds) == 0 {
					continue
				}

				goPackage = fds[0].AsFileDescriptorProto().GetOptions().GetGoPackage()

				command += fmt.Sprintf(" --go_opt=module=%s --go_opt=M%s=%s --go-grpc_opt=module=%s --go-grpc_opt=M%s=%s",
					jr.Module,
					rel,
					func() string {
						if strings.HasPrefix(goPackage, jr.Module) {
							return goPackage
						}
						return filepath.ToSlash(filepath.Join(jr.Module, "internal", goPackage))
					}(),
					jr.Module,
					rel,
					func() string {
						if strings.HasPrefix(goPackage, jr.Module) {
							return goPackage
						}
						return filepath.ToSlash(filepath.Join(jr.Module, "internal", goPackage))
					}())
			}

			if len(config.C.Gen.ProtoInclude) > 0 {
				command += fmt.Sprintf(" -I%s ", strings.Join(config.C.Gen.ProtoInclude, " -I"))
			}

			logx.Debug(command)

			_, err = execx.Run(command, config.C.Wd())
			if err != nil {
				return nil, err
			}
		}

		for _, file := range allServerFiles {
			if filepath.Clean(file.DescFilepath) == filepath.Clean(v) {
				if _, ok := genCodeProtoSpecMap[file.DescFilepath]; ok {
					if err = jr.removeServerSuffix(file.Path); err != nil {
						goctlconsole.Warning("[warning]: remove server suffix %s meet error %v", file.Path, err)
						continue
					}
				}
			}
		}

		for _, file := range allLogicFiles {
			if _, ok := genCodeProtoSpecMap[file.DescFilepath]; ok {
				if err := jr.removeLogicSuffix(file.Path); err != nil {
					goctlconsole.Warning("[warning]: remove logic suffix %s meet error %v", file.Path, err)
					continue
				}
			}
		}

		if lo.Contains(genCodeProtoFiles, v) {
			for _, file := range allLogicFiles {
				if err = jr.changeLogicTypes(file); err != nil {
					goctlconsole.Warning("[warning]: change logic types %s meet error %v", file.Path, err)
					continue
				}
			}
		}

		// # gen proto descriptor
		if lo.Contains(genCodeProtoFiles, v) {
			if jzerodesc.IsNeedGenProtoDescriptor(protoSpecMap[v]) {
				if !pathx.FileExists(jzerodesc.GetProtoDescriptorPath(v)) {
					_ = os.MkdirAll(filepath.Dir(jzerodesc.GetProtoDescriptorPath(v)), 0o755)
				}
				protocCommand := fmt.Sprintf("protoc --include_imports -I%s -I%s --descriptor_set_out=%s %s",
					config.C.ProtoDir(),
					filepath.Join(config.C.ProtoDir(), "third_party"),
					jzerodesc.GetProtoDescriptorPath(v),
					v,
				)
				if len(config.C.Gen.ProtoInclude) > 0 {
					protocCommand += fmt.Sprintf(" -I%s", strings.Join(config.C.Gen.ProtoInclude, " -I"))
				}
				logx.Debug(protocCommand)
				_, err = execx.Run(protocCommand, config.C.Wd())
				if err != nil {
					return nil, err
				}
			}
		}

		for _, s := range protoSpecMap[v].Service {
			serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, strings.ToLower(s.Name), jr.Module, strings.ToLower(s.Name)))
			registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%s(ctx))", filepath.Base(protoSpecMap[v].GoPackage), stringx.FirstUpper(s.Name), strings.ToLower(s.Name), stringx.FirstUpper(stringx.ToCamel(s.Name))))
		}
		pbImports = append(pbImports, fmt.Sprintf(`"%s/internal/%s"`, jr.Module, strings.TrimPrefix(protoSpecMap[v].GoPackage, "./")))
	}

	if len(genCodeProtoFiles) > 0 {
		if !config.C.Quiet {
			fmt.Println(console.Green("Done"))
		}
	}

	if pathx.FileExists(config.C.ProtoDir()) {
		if err = jr.genServer(serverImports, pbImports, registerServers); err != nil {
			return nil, err
		}
		// gen common proto pb
		if err = jr.genNoRpcServiceExcludeThirdPartyProto(config.C.ProtoDir()); err != nil {
			return nil, err
		}
		if err = jr.genApiMiddlewares(protoFiles); err != nil {
			return nil, err
		}
	}

	return protoSpecMap, nil
}
