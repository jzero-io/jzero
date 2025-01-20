package jparser

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"

	gconfig "github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser/api"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser/gateway"
	"github.com/jzero-io/jzero/internal/gen/gensdk/vars"
	"github.com/jzero-io/jzero/pkg/stringx"
)

func Parse(config *config.Config, fds []*desc.FileDescriptor, apiSpecs []*spec.ApiSpec) (vars.ScopeResourceHTTPInterfaceMap, error) {
	interfaces, err := genHTTPInterfaces(config, fds, apiSpecs)
	if err != nil {
		return nil, err
	}
	return convertToMap(interfaces), nil
}

func genHTTPInterfaces(config *config.Config, fds []*desc.FileDescriptor, apiSpecs []*spec.ApiSpec) ([]*vars.HTTPInterface, error) {
	var httpInterfaces []*vars.HTTPInterface

	for _, fd := range fds {
		services := fd.GetServices()
		for _, service := range services {
			methods := service.GetMethods()
			for _, method := range methods {
				ext := proto.GetExtension(method.GetMethodOptions(), annotations.E_Http)
				var httpInterface vars.HTTPInterface
				switch rule := ext.(type) {
				case *annotations.HttpRule:
					if rule == nil {
						continue
					}

					switch httpRule := rule.GetPattern().(type) {
					case *annotations.HttpRule_Get:
						httpInterface.Method = http.MethodGet
						httpInterface.URL = httpRule.Get
					case *annotations.HttpRule_Post:
						httpInterface.Method = http.MethodPost
						httpInterface.URL = httpRule.Post
					case *annotations.HttpRule_Put:
						httpInterface.Method = http.MethodPut
						httpInterface.URL = httpRule.Put
					case *annotations.HttpRule_Delete:
						httpInterface.Method = http.MethodDelete
						httpInterface.URL = httpRule.Delete
					case *annotations.HttpRule_Patch:
						httpInterface.Method = http.MethodPatch
						httpInterface.URL = httpRule.Patch
					}

					var requestBodyName string
					if (httpInterface.Method == http.MethodPost || httpInterface.Method == http.MethodPut || httpInterface.Method == http.MethodPatch) && rule.Body != "*" {
						for _, v := range method.GetInputType().GetFields() {
							if rule.Body == v.GetName() {
								requestBodyName = v.GetName()
							}
						}
					} else {
						requestBodyName = method.GetInputType().GetName()
					}
					httpInterface.IsStreamClient = method.IsClientStreaming()
					httpInterface.IsStreamServer = method.IsServerStreaming()
					httpInterface.Request = &vars.Request{
						RealBodyName: requestBodyName,
						Name:         stringx.FirstUpper(method.GetInputType().GetName()),
						Body:         rule.Body,
						Type:         "proto",
						Package:      strings.TrimPrefix(*service.GetFile().GetFileOptions().GoPackage, "./"),
						FullName:     fmt.Sprintf("param *%s.%s", filepath.Base(strings.TrimPrefix(*service.GetFile().GetFileOptions().GoPackage, "./")), stringx.FirstUpper(method.GetInputType().GetName())),
					}
					httpInterface.Response = &vars.Response{
						Package: strings.TrimPrefix(*service.GetFile().GetFileOptions().GoPackage, "./"),
					}
					httpInterface.Response.FullName = BuildProtoResponseFullName(httpInterface.Response.Package, stringx.FirstUpper(method.GetOutputType().GetName()))
				}
				httpInterface.MethodName = method.GetName()
				httpInterface.Scope = vars.Scope(gconfig.C.Gen.Sdk.Scope)
				httpInterface.Resource = vars.Resource(service.GetName())

				pathParams, err := gateway.PathParam(httpInterface.URL)
				if err != nil {
					return nil, nil
				}
				httpInterface.PathParams = pathParams
				queryParams := gateway.CreateQueryParams(method)
				httpInterface.QueryParams = queryParams

				httpInterfaces = append(httpInterfaces, &httpInterface)
			}
		}
	}

	for _, apiSpec := range apiSpecs {
		for _, group := range apiSpec.Service.Groups {
			for _, route := range group.Routes {
				path, _ := url.JoinPath(group.Annotation.Properties["prefix"], route.Path)

				resource := vars.Resource(stringx.ToCamel(group.Annotation.Properties["group"]))
				if resource == "" {
					resource = "api"
				}
				httpInterface := vars.HTTPInterface{
					Scope:      vars.Scope(gconfig.C.Gen.Sdk.Scope),
					Resource:   resource,
					Method:     strings.ToUpper(route.Method),
					URL:        path,
					MethodName: strings.TrimSuffix(route.Handler, "Handler"),
					Comments:   route.AtDoc.Text,
				}
				if route.RequestType != nil {
					httpInterface.Request = &vars.Request{
						Name:         stringx.FirstUpper(route.RequestType.Name()),
						RealBodyName: stringx.FirstUpper(route.RequestType.Name()),
						Package:      "types",
						Type:         "api",
						FullName:     fmt.Sprintf("param %s.%s", "types", stringx.FirstUpper(route.RequestType.Name())),
					}
					if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok && goPackage != "" {
						httpInterface.Request.Package = goPackage
						httpInterface.Request.FullName = fmt.Sprintf("param %s.%s", strings.ToLower(strings.ReplaceAll(apiSpec.Info.Properties["go_package"], "/", "")), stringx.FirstUpper(route.RequestType.Name()))
					}

					if strings.ToUpper(route.Method) == http.MethodPost || strings.ToUpper(route.Method) == http.MethodPut || strings.ToUpper(route.Method) == http.MethodPatch {
						httpInterface.Request.Body = "*"
					}
				} else {
					httpInterface.IsStreamClient = true
					httpInterface.Request = &vars.Request{}
				}

				if route.ResponseType != nil {
					httpInterface.Response = &vars.Response{
						FullName: BuildApiResponseFullName(apiSpec, route.ResponseType, false),
						Package:  "types",
					}
					if goPackage, ok := apiSpec.Info.Properties["go_package"]; ok && goPackage != "" {
						httpInterface.Response.Package = goPackage
						httpInterface.Response.FullName = BuildApiResponseFullName(apiSpec, route.ResponseType, true)
					}
				} else {
					httpInterface.IsStreamServer = true
				}

				pathParams, err := api.CreatePathParam(httpInterface.URL, &route)
				if err != nil {
					return nil, nil
				}
				httpInterface.PathParams = pathParams
				queryParams := api.CreateQueryParams(&route)
				httpInterface.QueryParams = queryParams

				httpInterfaces = append(httpInterfaces, &httpInterface)
			}
		}
	}
	return httpInterfaces, nil
}

func BuildProtoResponseFullName(goPackage, responseTypeName string) string {
	return fmt.Sprintf("*%s.%s", filepath.Base(goPackage), responseTypeName)
}

func BuildApiResponseFullName(apiSpec *spec.ApiSpec, t spec.Type, splitApiTypesDir bool) string {
	switch v := t.(type) {
	case spec.PrimitiveType:
		return "*" + t.Name()
	case spec.ArrayType:
		if _, ok := v.Value.(spec.PrimitiveType); ok {
			return t.Name()
		}
		if splitApiTypesDir {
			return fmt.Sprintf("[]*%s.%s", strings.ToLower(strings.ReplaceAll(apiSpec.Info.Properties["go_package"], "/", "")), strings.TrimPrefix(t.Name(), "[]"))
		}
		return fmt.Sprintf("[]*types.%s", strings.TrimPrefix(t.Name(), "[]"))
	default:
		if splitApiTypesDir {
			return fmt.Sprintf("*%s.%s", strings.ToLower(strings.ReplaceAll(apiSpec.Info.Properties["go_package"], "/", "")), stringx.FirstUpper(t.Name()))
		}
		return fmt.Sprintf("*types.%s", stringx.FirstUpper(strings.TrimPrefix(t.Name(), "*")))
	}
}

func convertToMap(interfaces []*vars.HTTPInterface) vars.ScopeResourceHTTPInterfaceMap {
	scopeResourceHTTPInterfaceMap := make(vars.ScopeResourceHTTPInterfaceMap)

	for _, inf := range interfaces {
		scope := inf.Scope
		resource := inf.Resource

		if _, ok := scopeResourceHTTPInterfaceMap[scope]; !ok {
			scopeResourceHTTPInterfaceMap[scope] = make(map[vars.Resource][]*vars.HTTPInterface)
		}

		scopeResourceHTTPInterfaceMap[scope][resource] = append(scopeResourceHTTPInterfaceMap[scope][resource], inf)
	}

	return scopeResourceHTTPInterfaceMap
}
