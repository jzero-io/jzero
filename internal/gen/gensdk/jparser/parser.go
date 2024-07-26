package jparser

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser/api"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser/gateway"
	"github.com/jzero-io/jzero/internal/gen/gensdk/vars"

	"github.com/jhump/protoreflect/desc"
	"github.com/jzero-io/jzero/pkg/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
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
					}

					var requestBodyName string
					if (httpInterface.Method == http.MethodPost) && rule.Body != "*" {
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
						// param *{{.Request.Package | base}}.{{.Request.Name}}
						FullName: fmt.Sprintf("param *%s.%s", filepath.Base(strings.TrimPrefix(*service.GetFile().GetFileOptions().GoPackage, "./")), stringx.FirstUpper(method.GetInputType().GetName())),
					}
					httpInterface.Response = &vars.Response{
						Package: strings.TrimPrefix(*service.GetFile().GetFileOptions().GoPackage, "./"),
					}
					httpInterface.Response.FullName = BuildProtoResponseFullName(httpInterface.Response.Package, stringx.FirstUpper(method.GetOutputType().GetName()))
					httpInterface.Response.FakeFullName = BuildProtoFakeFullName(httpInterface.Response.Package, stringx.FirstUpper(method.GetOutputType().GetName()))
					httpInterface.Response.FakeReturnName = BuildProtoFakeReturnName(service.GetName(), stringx.FirstUpper(method.GetOutputType().GetName()))
				}
				httpInterface.MethodName = method.GetName()
				httpInterface.Scope = vars.Scope(config.Scope)
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

				httpInterface := vars.HTTPInterface{
					Scope:      vars.Scope(config.Scope),
					Resource:   vars.Resource(stringx.ToCamel(group.Annotation.Properties["group"])),
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
					if strings.ToUpper(route.Method) == http.MethodPost {
						httpInterface.Request.Body = "*"
					}
				} else {
					httpInterface.IsStreamClient = true
					httpInterface.Request = &vars.Request{}
				}

				if route.ResponseType != nil {
					httpInterface.Response = &vars.Response{
						FakeReturnName: BuildApiFakeReturnName(group.GetAnnotation("group"), httpInterface.MethodName, route.ResponseType),
						FakeFullName:   BuildApiFakeFullName(route.ResponseType),
						FullName:       BuildApiResponseFullName(route.ResponseType),
						Package:        "types",
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

func BuildProtoResponseFullName(goPackage string, responseTypeName string) string {
	return fmt.Sprintf("*%s.%s", filepath.Base(goPackage), responseTypeName)
}

func BuildProtoFakeFullName(goPackage string, responseTypeName string) string {
	return fmt.Sprintf("&%s.%s", filepath.Base(goPackage), responseTypeName)
}

func BuildApiResponseFullName(t spec.Type) string {
	switch v := t.(type) {
	case spec.PrimitiveType:
		return "*" + t.Name()
	case spec.ArrayType:
		if _, ok := v.Value.(spec.PrimitiveType); ok {
			return t.Name()
		}
		return fmt.Sprintf("[]*types.%s", strings.TrimPrefix(t.Name(), "[]"))
	default:
		return fmt.Sprintf("*types.%s", stringx.FirstUpper(strings.TrimPrefix(t.Name(), "*")))
	}
}

func BuildProtoFakeReturnName(serviceName string, responseTypeName string) string {
	return fmt.Sprintf("FakeReturn%s%s", stringx.FirstUpper(serviceName), stringx.FirstUpper(responseTypeName))
}

func BuildApiFakeReturnName(group string, method string, t spec.Type) string {
	return fmt.Sprintf("FakeReturn%s%s%s", stringx.FirstUpper(stringx.ToCamel(group)), stringx.FirstUpper(method), stringx.FirstUpper(t.Name()))
}

func BuildApiFakeFullName(t spec.Type) string {
	switch v := t.(type) {
	case spec.PrimitiveType:
		return "*" + t.Name()
	case spec.ArrayType:
		if value, ok := v.Value.(spec.PrimitiveType); ok {
			return value.RawName
		}
		return fmt.Sprintf("[]*types.%s", strings.TrimPrefix(t.Name(), "[]"))
	default:
		return fmt.Sprintf("&types.%s", stringx.FirstUpper(strings.TrimPrefix(t.Name(), "*")))
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
