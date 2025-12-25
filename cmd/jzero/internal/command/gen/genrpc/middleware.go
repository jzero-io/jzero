package genrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/orderedmap"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	jzeroapi "github.com/jzero-io/desc/proto/api"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"google.golang.org/protobuf/proto"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type JzeroProtoApiMiddleware struct {
	Name   string
	Routes []string
}

func (jr *JzeroRpc) genApiMiddlewares(protoFiles []string) (err error) {
	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser

	protoParser.InferImportPaths = false

	var files []string
	for _, protoFilename := range protoFiles {
		rel, err := filepath.Rel(filepath.Join("desc", "proto"), protoFilename)
		if err != nil {
			return err
		}
		files = append(files, rel)
	}

	protoParser.ImportPaths = []string{filepath.Join("desc", "proto"), filepath.Join("desc", "proto", "third_party")}
	if len(config.C.Gen.ProtoInclude) > 0 {
		protoParser.ImportPaths = append(protoParser.ImportPaths, config.C.Gen.ProtoInclude...)
	}
	protoParser.IncludeSourceCodeInfo = true
	fds, err = protoParser.ParseFiles(files...)
	if err != nil {
		return err
	}

	var httpMiddlewares []JzeroProtoApiMiddleware
	var zrpcMiddlewares []JzeroProtoApiMiddleware

	httpMapMiddlewares := orderedmap.New()
	zrpcMapMiddlewares := orderedmap.New()

	for _, fd := range fds {
		descriptorProto := fd.AsFileDescriptorProto()

		var methodUrls []string
		var fullMethods []string

		for _, service := range descriptorProto.GetService() {
			for _, method := range service.GetMethod() {
				methodUrls = append(methodUrls, jzerodesc.GetRpcMethodUrl(method))
				fullMethods = append(fullMethods, fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), service.GetName(), method.GetName()))

				httpGroupExt := proto.GetExtension(service.GetOptions(), jzeroapi.E_HttpGroup)
				switch rule := httpGroupExt.(type) {
				case *jzeroapi.HttpRule:
					if rule != nil {
						split := strings.Split(strings.ReplaceAll(rule.Middleware, " ", ""), ",")
						for _, m := range split {
							if urls, ok := httpMapMiddlewares.Get(m); ok {
								urls = append(urls.([]string), methodUrls...)
								httpMapMiddlewares.Set(m, urls)
							} else {
								httpMapMiddlewares.Set(m, methodUrls)
							}
						}
					}
				}

				zrpcGroupExt := proto.GetExtension(service.GetOptions(), jzeroapi.E_ZrpcGroup)
				switch rule := zrpcGroupExt.(type) {
				case *jzeroapi.ZrpcRule:
					if rule != nil {
						split := strings.Split(strings.ReplaceAll(rule.Middleware, " ", ""), ",")
						for _, m := range split {
							if fms, ok := zrpcMapMiddlewares.Get(m); ok {
								fms = append(fms.([]string), fullMethods...)
								zrpcMapMiddlewares.Set(m, fms)
							} else {
								zrpcMapMiddlewares.Set(m, fullMethods)
							}
						}
					}
				}

				httpExt := proto.GetExtension(method.GetOptions(), jzeroapi.E_Http)
				switch rule := httpExt.(type) {
				case *jzeroapi.HttpRule:
					if rule != nil {
						split := strings.Split(strings.ReplaceAll(rule.Middleware, " ", ""), ",")
						for _, m := range split {
							if urls, ok := httpMapMiddlewares.Get(m); ok {
								urls = append(urls.([]string), jzerodesc.GetRpcMethodUrl(method))
								httpMapMiddlewares.Set(m, urls)
							} else {
								httpMapMiddlewares.Set(m, []string{jzerodesc.GetRpcMethodUrl(method)})
							}
						}
					}
				}

				zrpcExt := proto.GetExtension(method.GetOptions(), jzeroapi.E_Zrpc)
				switch rule := zrpcExt.(type) {
				case *jzeroapi.ZrpcRule:
					if rule != nil {
						split := strings.Split(strings.ReplaceAll(rule.Middleware, " ", ""), ",")
						for _, m := range split {
							if urls, ok := zrpcMapMiddlewares.Get(m); ok {
								urls = append(urls.([]string), fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), service.GetName(), method.GetName()))
								zrpcMapMiddlewares.Set(m, urls)
							} else {
								zrpcMapMiddlewares.Set(m, []string{fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), service.GetName(), method.GetName())})
							}
						}
					}
				}
			}
		}
	}

	// order and unique and transfer to httpMiddlewares and zrpcMiddlewares
	httpMiddlewares = processMiddlewares(httpMapMiddlewares)
	zrpcMiddlewares = processMiddlewares(zrpcMapMiddlewares)

	if len(httpMiddlewares) == 0 && len(zrpcMiddlewares) == 0 {
		return nil
	}

	if !config.C.Quiet {
		fmt.Printf("%s to generate internal/middleware/middleware_gen.go\n", console.Green("Start"))
	}

	for _, v := range httpMiddlewares {
		template, err := templatex.ParseTemplate(filepath.Join("rpc", "middleware_http.go.tpl"), map[string]any{
			"Name": v.Name,
		}, embeded.ReadTemplateFile(filepath.Join("rpc", "middleware_http.go.tpl")))
		if err != nil {
			return err
		}
		process, err := gosimports.Process("", template, nil)
		if err != nil {
			return err
		}
		namingFormat, _ := format.FileNamingFormat(config.C.Style, v.Name+"Middleware")
		if !pathx.FileExists(filepath.Join("internal", "middleware", namingFormat+".go")) {
			err = os.WriteFile(filepath.Join("internal", "middleware", namingFormat+".go"), process, 0o644)
			if err != nil {
				return err
			}
		}
	}

	for _, v := range zrpcMiddlewares {
		template, err := templatex.ParseTemplate(filepath.Join("rpc", "middleware_zrpc.go.tpl"), map[string]any{
			"Name": v.Name,
		}, embeded.ReadTemplateFile(filepath.Join("rpc", "middleware_zrpc.go.tpl")))
		if err != nil {
			return err
		}

		process, err := gosimports.Process("", template, nil)
		if err != nil {
			return err
		}
		namingFormat, _ := format.FileNamingFormat(config.C.Style, v.Name+"Middleware")
		if !pathx.FileExists(filepath.Join("internal", "middleware", namingFormat+".go")) {
			err = os.WriteFile(filepath.Join("internal", "middleware", namingFormat+".go"), process, 0o644)
			if err != nil {
				return err
			}
		}
	}

	template, err := templatex.ParseTemplate(filepath.Join("rpc", "middleware_gen.go.tpl"), map[string]any{
		"HttpMiddlewares": httpMiddlewares,
		"ZrpcMiddlewares": zrpcMiddlewares,
	}, embeded.ReadTemplateFile(filepath.Join("rpc", "middleware_gen.go.tpl")))
	if err != nil {
		return err
	}

	process, err := gosimports.Process("", template, nil)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join("internal", "middleware", "middleware_gen.go"), process, 0o644)
	if err != nil {
		return err
	}

	if !config.C.Quiet {
		fmt.Printf("%s\n", console.Green("Done"))
	}
	return nil
}

func processMiddlewares(middlewareMap *orderedmap.OrderedMap) []JzeroProtoApiMiddleware {
	var result []JzeroProtoApiMiddleware

	for _, m := range middlewareMap.Keys() {
		v, _ := middlewareMap.Get(m)
		v = lo.Uniq(v.([]string))
		result = append(result, JzeroProtoApiMiddleware{Name: m, Routes: v.([]string)})
	}
	return result
}
