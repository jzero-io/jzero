package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jzero-io/jzero/config"
	"github.com/rinchsan/gosimports"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	jzeroapi "github.com/jzero-io/desc/proto/jzero/api"

	"google.golang.org/protobuf/proto"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
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
	Wd                 string
	Module             string
	Style              string
	RemoveSuffix       bool
	ChangeReplaceTypes bool
}

type JzeroProtoApiMiddleware struct {
	Name   string
	Routes []string
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
	var allServerFiles []ServerFile
	var allLogicFiles []LogicFile

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

		if jr.ChangeReplaceTypes {
			for _, file := range allLogicFiles {
				if err := jr.changeReplaceLogicGoTypes(file); err != nil {
					continue
				}
			}
		}

		// # gen proto descriptor
		if isNeedGenProtoDescriptor(parse) {
			if !pathx.FileExists(generateProtoDescriptorPath(v)) {
				_ = os.MkdirAll(filepath.Dir(generateProtoDescriptorPath(v)), 0o755)
			}
			protocCommand := fmt.Sprintf("protoc --include_imports -I%s --descriptor_set_out=%s %s",
				protoDirPath,
				generateProtoDescriptorPath(v),
				v,
			)
			_, err = execx.Run(protocCommand, jr.Wd)
			if err != nil {
				return err
			}
		}

		for _, s := range parse.Service {
			serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/internal/server/%s"`, strings.ToLower(s.Name), jr.Module, strings.ToLower(s.Name)))
			if jr.RemoveSuffix {
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%s(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), strings.ToLower(s.Name), stringx.FirstUpper(s.Name)))
			} else {
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), strings.ToLower(s.Name), stringx.FirstUpper(s.Name)))
			}
		}
		pbImports = append(pbImports, fmt.Sprintf(`"%s/internal/%s"`, jr.Module, strings.TrimPrefix(parse.GoPackage, "./")))
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

func generateProtoDescriptorPath(protoPath string) string {
	rel, err := filepath.Rel(filepath.Join("desc", "proto"), protoPath)
	if err != nil {
		return ""
	}

	return filepath.Join("desc", "pb", strings.TrimSuffix(rel, ".proto")+".pb")
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
				Path:             fp,
				Handler:          rpc.Name,
				Group:            service.Name,
				ClientStream:     rpc.StreamsRequest,
				ServerStream:     rpc.StreamsReturns,
				ResponseTypeName: rpc.ReturnsType,
				RequestTypeName:  rpc.RequestType,
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

func (jr *JzeroRpc) changeReplaceLogicGoTypes(file LogicFile) error {
	fp := file.Path // logic file path
	if jr.RemoveSuffix {
		// Get the new file name of the file (without the 5 characters(Logic or logic) before the ".go" extension)
		fp = file.Path[:len(file.Path)-8]
		// patch
		fp = strings.TrimSuffix(fp, "_")
		fp = strings.TrimSuffix(fp, "-")
		fp += ".go"
	}

	fset := token.NewFileSet()

	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	var needModify bool
	ast.Inspect(f, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncDecl); ok && fn.Recv != nil {
			if fn.Name.Name == file.Handler {
				if fn.Type != nil && fn.Type.Params != nil {
					for _, param := range fn.Type.Params.List {
						if starExpr, ok := param.Type.(*ast.StarExpr); ok {
							if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
								if selectorExpr.Sel.Name != file.RequestTypeName {
									selectorExpr.Sel.Name = file.RequestTypeName
									needModify = true
								}
							}
						}
					}
				}

				if fn.Type != nil && fn.Type.Results != nil {
					for _, result := range fn.Type.Results.List {
						if starExpr, ok := result.Type.(*ast.StarExpr); ok {
							if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
								if selectorExpr.Sel.Name != file.ResponseTypeName {
									selectorExpr.Sel.Name = file.ResponseTypeName
									needModify = true
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	if needModify {
		// Write the modified AST back to the file
		buf := bytes.NewBuffer(nil)
		if err := goformat.Node(buf, fset, f); err != nil {
			return err
		}

		if err = os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (jr *JzeroRpc) genApiMiddlewares(protoFilenames []string) (err error) {
	// first: parse all proto, find option
	fmt.Printf("%s to generate internal/middleware/middleware_gen.go\n", color.WithColor("Start", color.FgGreen))

	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser

	protoParser.InferImportPaths = false

	var files []string
	for _, protoFilename := range protoFilenames {
		rel, err := filepath.Rel(filepath.Join("desc", "proto"), protoFilename)
		if err != nil {
			return err
		}
		files = append(files, rel)
	}

	protoParser.ImportPaths = []string{filepath.Join("desc", "proto")}
	protoParser.IncludeSourceCodeInfo = true
	fds, err = protoParser.ParseFiles(files...)
	if err != nil {
		return err
	}

	var httpMiddlewares []JzeroProtoApiMiddleware
	var zrpcMiddlewares []JzeroProtoApiMiddleware

	httpMapMiddlewares := make(map[string][]string)
	zrpcMapMiddlewares := make(map[string][]string)

	for _, fd := range fds {
		descriptorProto := fd.AsFileDescriptorProto()

		var methodUrls []string
		var fullMethods []string

		for _, service := range descriptorProto.GetService() {
			for _, method := range service.GetMethod() {
				methodUrls = append(methodUrls, getRpcMethodUrl(method))
				fullMethods = append(fullMethods, fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), service.GetName(), method.GetName()))

				httpExt := proto.GetExtension(method.GetOptions(), jzeroapi.E_Http)
				switch rule := httpExt.(type) {
				case *jzeroapi.HttpRule:
					if rule != nil {
						httpMapMiddlewares[rule.Middleware] = append(httpMapMiddlewares[rule.Middleware], getRpcMethodUrl(method))
					}
				}
				zrpcExt := proto.GetExtension(method.GetOptions(), jzeroapi.E_Zrpc)
				switch rule := zrpcExt.(type) {
				case *jzeroapi.ZrpcRule:
					if rule != nil {
						zrpcMapMiddlewares[rule.Middleware] = append(zrpcMapMiddlewares[rule.Middleware], fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), service.GetName(), method.GetName()))
					}
				}
			}
			httpGroupExt := proto.GetExtension(service.GetOptions(), jzeroapi.E_HttpGroup)
			switch rule := httpGroupExt.(type) {
			case *jzeroapi.HttpRule:
				if rule != nil {
					httpMapMiddlewares[rule.Middleware] = append(httpMapMiddlewares[rule.Middleware], methodUrls...)
				}
			}

			zrpcGroupExt := proto.GetExtension(service.GetOptions(), jzeroapi.E_ZrpcGroup)
			switch rule := zrpcGroupExt.(type) {
			case *jzeroapi.ZrpcRule:
				if rule != nil {
					zrpcMapMiddlewares[rule.Middleware] = append(zrpcMapMiddlewares[rule.Middleware], fullMethods...)
				}
			}
		}
	}

	// sort and unique and transfer to httpMiddlewares and zrpcMiddlewares
	httpMiddlewares = processMiddlewares(httpMapMiddlewares)
	zrpcMiddlewares = processMiddlewares(zrpcMapMiddlewares)

	for _, v := range httpMiddlewares {
		template, err := templatex.ParseTemplate(map[string]interface{}{
			"Name": v.Name,
		}, embeded.ReadTemplateFile(filepath.Join("plugins", "api", "middleware_http.go.tpl")))
		if err != nil {
			return err
		}

		process, err := gosimports.Process("", template, &gosimports.Options{
			Comments:   true,
			FormatOnly: true,
		})
		if err != nil {
			return err
		}
		namingFormat, _ := format.FileNamingFormat(config.C.Gen.Style, v.Name+"Middleware")
		if !pathx.FileExists(filepath.Join("internal", "middleware", namingFormat+".go")) {
			err = os.WriteFile(filepath.Join("internal", "middleware", namingFormat+".go"), process, 0o644)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range zrpcMiddlewares {
		template, err := templatex.ParseTemplate(map[string]interface{}{
			"Name": v.Name,
		}, embeded.ReadTemplateFile(filepath.Join("plugins", "api", "middleware_zrpc.go.tpl")))
		if err != nil {
			return err
		}

		process, err := gosimports.Process("", template, &gosimports.Options{
			Comments:   true,
			FormatOnly: true,
		})
		if err != nil {
			return err
		}
		namingFormat, _ := format.FileNamingFormat(config.C.Gen.Style, v.Name+"Middleware")
		if !pathx.FileExists(filepath.Join("internal", "middleware", namingFormat+".go")) {
			err = os.WriteFile(filepath.Join("internal", "middleware", namingFormat+".go"), process, 0o644)
			if err != nil {
				return err
			}
		}
	}

	template, err := templatex.ParseTemplate(map[string]interface{}{
		"HttpMiddlewares": httpMiddlewares,
		"ZrpcMiddlewares": zrpcMiddlewares,
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "api", "middleware_gen.go.tpl")))
	if err != nil {
		return err
	}

	process, err := gosimports.Process("", template, &gosimports.Options{
		Comments:   true,
		FormatOnly: true,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join("internal", "middleware", "middleware_gen.go"), process, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func processMiddlewares(middlewareMap map[string][]string) []JzeroProtoApiMiddleware {
	var result []JzeroProtoApiMiddleware

	for name, routes := range middlewareMap {
		uniqueRoutes := uniqueAndSort(routes)
		result = append(result, JzeroProtoApiMiddleware{Name: name, Routes: uniqueRoutes})
	}

	// Sort the middleware list by name if needed
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func uniqueAndSort(input []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, item := range input {
		uniqueMap[item] = struct{}{}
	}

	uniqueList := make([]string, 0, len(uniqueMap))
	for item := range uniqueMap {
		uniqueList = append(uniqueList, item)
	}

	sort.Strings(uniqueList)
	return uniqueList
}
