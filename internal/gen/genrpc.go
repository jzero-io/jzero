package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"gopkg.in/op/go-logging.v1"

	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/jzero-io/jzero/pkg/templatex"
	yq "github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type (
	ImportLines   []string
	RegisterLines []string
)

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

type ServerFile struct {
	Path string
}

type JzeroRpc struct {
	Wd           string
	Module       string
	Style        string
	RemoveSuffix bool
	Etc          string // 配置文件路径
}

func (jr *JzeroRpc) Gen() error {
	protoDirPath := filepath.Join("desc", "proto")
	protoFilenames, err := GetProtoFilepath(protoDirPath)
	if err != nil {
		return err
	}

	var serverImports ImportLines
	var pbImports ImportLines
	var registerServers RegisterLines
	var protoDescriptorPaths []string

	var allServerFiles []ServerFile
	var allLogicFiles []LogicFile

	for _, v := range protoFilenames {
		// parse proto
		protoParser := rpcparser.NewDefaultProtoParser()
		parse, err := protoParser.Parse(v, true)
		if err != nil {
			continue
		}

		allLogicFiles, err = jr.GetAllLogicFiles(parse)
		if err != nil {
			return err
		}

		allServerFiles, err = jr.GetAllServerFiles(parse)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate proto code. \n%s proto file %s\n", color.WithColor("Start", color.FgGreen), color.WithColor("Using", color.FgGreen), v)
		zrpcOut := "."

		command := fmt.Sprintf("goctl rpc protoc %s -I%s --go_out=%s --go-grpc_out=%s --zrpc_out=%s --client=false --home %s -m --style %s ",
			v,
			protoDirPath,
			filepath.Join("internal"),
			filepath.Join("internal"),
			zrpcOut,
			filepath.Join(embeded.Home, "go-zero"),
			jr.Style)

		fileBase := v[0 : len(v)-len(path.Ext(v))]

		_, err = execx.Run(command, jr.Wd)
		if err != nil {
			return err
		}

		command = fmt.Sprintf("protoc %s -I%s --validate_out=%s",
			v,
			protoDirPath,
			"lang=go:internal",
		)
		_, err = execx.Run(command, jr.Wd)
		if err != nil {
			return err
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))

		if jr.RemoveSuffix {
			for _, file := range allServerFiles {
				if err := jr.rewriteServerGo(file.Path); err != nil {
					continue
				}
			}
			for _, file := range allLogicFiles {
				if err := jr.rewriteLogicGo(file.Path); err != nil {
					continue
				}
			}
		}

		// # gen proto descriptor
		if isNeedGenProtoDescriptor(parse) {
			protocCommand := fmt.Sprintf("protoc --include_imports -I%s --descriptor_set_out=%s.pb %s",
				protoDirPath,
				fileBase,
				v,
			)
			protoDescriptorPaths = append(protoDescriptorPaths, fmt.Sprintf("%s.pb", strings.TrimSuffix(v, ".proto")))
			_, err = execx.Run(protocCommand, jr.Wd)
			if err != nil {
				return err
			}

			// update gateway upstream protosets
			err = jr.updateGatewayUpstreams(protoDescriptorPaths)
			if err != nil {
				console.Warning("[warning] update gateway upstreams meet error: %v", err)
			}
		}

		for _, s := range parse.Service {
			serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, s.Name, jr.Module, s.Name))
			if jr.RemoveSuffix {
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%s(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), s.Name, stringx.FirstUpper(s.Name)))
			} else {
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), s.Name, stringx.FirstUpper(s.Name)))
			}
		}
		pbImports = append(pbImports, fmt.Sprintf(`"%s/internal/%s"`, jr.Module, strings.TrimPrefix(parse.GoPackage, "./")))
	}

	if pathx.FileExists(protoDirPath) {
		if err = jr.genServer(serverImports, pbImports, registerServers); err != nil {
			return err
		}
	}
	return nil
}

func (jr *JzeroRpc) genServer(serverImports ImportLines, pbImports ImportLines, registerServers RegisterLines) error {
	fmt.Printf("%s to generate internal/server/server.go\n", color.WithColor("Start", color.FgGreen))
	serverFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module":          jr.Module,
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("app", "internal", "server", "server.go.tpl")))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(jr.Wd, "internal", "server", "server.go"), serverFile, 0o644)
	if err != nil {
		return err
	}
	fmt.Printf("%s", color.WithColor("Done\n", color.FgGreen))
	return nil
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

func (jr *JzeroRpc) GetAllServerFiles(protoSpec rpcparser.Proto) ([]ServerFile, error) {
	var serverFiles []ServerFile
	for _, service := range protoSpec.Service {
		namingFormat, err := format.FileNamingFormat(jr.Style, service.Name+"Server")
		if err != nil {
			return nil, err
		}
		fp := filepath.Join(jr.Wd, "internal", "server", strings.ToLower(service.Name), namingFormat+".go")

		f := ServerFile{
			Path: fp,
		}

		serverFiles = append(serverFiles, f)
	}
	return serverFiles, nil
}

func (jr *JzeroRpc) GetAllLogicFiles(protoSpec rpcparser.Proto) ([]LogicFile, error) {
	var logicFiles []LogicFile
	for _, service := range protoSpec.Service {
		for _, rpc := range service.RPC {
			namingFormat, err := format.FileNamingFormat(jr.Style, rpc.Name+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(jr.Wd, "internal", "logic", strings.ToLower(service.Name), namingFormat+".go")

			f := LogicFile{
				Path:    fp,
				Handler: rpc.Name,
				Group:   service.Name,
			}

			logicFiles = append(logicFiles, f)
		}
	}
	return logicFiles, nil
}

func (jr *JzeroRpc) rewriteLogicGo(fp string) error {
	// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
	newFilePath := fp[:len(fp)-8]
	// patch
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	newFilePath += ".go"

	if pathx.FileExists(newFilePath) {
		_ = os.Remove(fp)
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasSuffix(fn.Name.Name, "Logic") {
			fn.Name.Name = strings.TrimSuffix(fn.Name.Name, "Logic")
			for _, result := range fn.Type.Results.List {
				if starExpr, ok := result.Type.(*ast.StarExpr); ok {
					if indent, ok := starExpr.X.(*ast.Ident); ok {
						indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Logic"))
					}
				}
			}
			for _, body := range fn.Body.List {
				if returnStmt, ok := body.(*ast.ReturnStmt); ok {
					for _, result := range returnStmt.Results {
						if unaryExpr, ok := result.(*ast.UnaryExpr); ok {
							if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								if indent, ok := compositeLit.Type.(*ast.Ident); ok {
									indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Logic"))
								}
							}
						}
					}
				}
			}
			return false
		}
		return true
	})

	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.GenDecl); ok && fn.Tok == token.TYPE {
			for _, s := range fn.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					ts.Name.Name = strings.TrimSuffix(ts.Name.Name, "Logic")
				}
			}
		}
		return true
	})

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil {
			for _, list := range fn.Recv.List {
				if starExpr, ok := list.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						ident.Name = util.Title(strings.TrimSuffix(ident.Name, "Logic"))
					}
				}
			}
		}
		return true
	})

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	return os.Rename(fp, newFilePath)
}

func (jr *JzeroRpc) rewriteServerGo(fp string) error {
	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && strings.HasSuffix(fn.Name.Name, "Server") {
			fn.Name.Name = strings.TrimSuffix(fn.Name.Name, "Server")
			for _, result := range fn.Type.Results.List {
				if starExpr, ok := result.Type.(*ast.StarExpr); ok {
					if indent, ok := starExpr.X.(*ast.Ident); ok {
						indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Server"))
					}
				}
			}
			for _, body := range fn.Body.List {
				if returnStmt, ok := body.(*ast.ReturnStmt); ok {
					for _, result := range returnStmt.Results {
						if unaryExpr, ok := result.(*ast.UnaryExpr); ok {
							if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								if indent, ok := compositeLit.Type.(*ast.Ident); ok {
									indent.Name = util.Title(strings.TrimSuffix(indent.Name, "Server"))
								}
							}
						}
					}
				}
			}
			return false
		}
		return true
	})

	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.GenDecl); ok && fn.Tok == token.TYPE {
			for _, s := range fn.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					ts.Name.Name = strings.TrimSuffix(ts.Name.Name, "Server")
				}
			}
		}
		return true
	})

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Recv != nil {
			for _, list := range fn.Recv.List {
				if starExpr, ok := list.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						ident.Name = util.Title(strings.TrimSuffix(ident.Name, "Server"))
					}
				}
			}
			// find handlerFunc body: l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
			for _, body := range fn.Body.List {
				if assignStmt, ok := body.(*ast.AssignStmt); ok {
					for _, rhs := range assignStmt.Rhs {
						if callExpr, ok := rhs.(*ast.CallExpr); ok {
							if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								selectorExpr.Sel.Name = strings.TrimSuffix(selectorExpr.Sel.Name, "Logic")
							}
						}
					}
				}
			}
		}
		return true
	})

	// Write the modified AST back to the file
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}

	// Get the new file name of the file (without the 5 characters(Server or server) before the ".go" extension)
	newFilePath := fp[:len(fp)-9]
	// patch
	newFilePath = strings.TrimSuffix(newFilePath, "_")
	newFilePath = strings.TrimSuffix(newFilePath, "-")
	if err = os.Rename(fp, newFilePath+".go"); err != nil {
		return err
	}

	return nil
}

func (jr *JzeroRpc) updateGatewayUpstreams(protoDescriptorPaths []string) error {
	logging.SetLevel(logging.CRITICAL, "")
	dec := yq.NewYamlDecoder(yq.NewDefaultYamlPreferences())

	file, err := os.ReadFile(jr.Etc)
	if err != nil {
		return err
	}

	if err = dec.Init(strings.NewReader(string(file))); err != nil {
		return err
	}

	node, err := dec.Decode()
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(protoDescriptorPaths)
	if err != nil {
		return err
	}

	result, _ := yq.NewAllAtOnceEvaluator().EvaluateNodes(fmt.Sprintf(`.Gateway.Upstreams.0.ProtoSets=%s`, string(marshal)), node)
	encoder := yq.NewYamlEncoder(yq.NewDefaultYamlPreferences())
	out := new(bytes.Buffer)
	printer := yq.NewPrinter(encoder, yq.NewSinglePrinterWriter(out))
	if err := printer.PrintResults(result); err != nil {
		return err
	}
	return os.WriteFile(jr.Etc, out.Bytes(), 0o644)
}
