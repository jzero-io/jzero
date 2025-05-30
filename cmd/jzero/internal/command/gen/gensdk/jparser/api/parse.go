package api

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/vars"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
)

func CreatePathParam(pattern string, route *spec.Route) ([]*vars.PathParam, error) {
	pathSegments := strings.Split(pattern, "/")
	params := make([]*vars.PathParam, 0)

	for i, segment := range pathSegments {
		if strings.HasPrefix(segment, ":") {
			param := &vars.PathParam{
				Index: i,
				Name:  segment[1:],
			}
			// get GoName
			if _, ok := route.RequestType.(spec.DefineStruct); ok {
				members := route.RequestType.(spec.DefineStruct).GetTagMembers("path")
				for _, member := range members {
					name, _ := member.GetPropertyName()
					if name == segment[1:] {
						param.GoName = stringx.FirstUpper(member.Name)
					}
				}
				params = append(params, param)
			}
		}
	}
	return params, nil
}

func CreateQueryParams(route *spec.Route) []*vars.QueryParam {
	if route.ResponseType == nil {
		return nil
	}

	if _, ok := route.RequestType.(spec.DefineStruct); !ok {
		return nil
	}

	return extractQueryParams(route.RequestType.(spec.DefineStruct))
}

func extractQueryParams(defineStruct spec.DefineStruct) []*vars.QueryParam {
	members := defineStruct.GetTagMembers("form")
	params := make([]*vars.QueryParam, 0)
	for _, member := range members {
		name, _ := member.GetPropertyName()
		if name != "" {
			param := &vars.QueryParam{
				GoName: stringx.FirstUpper(member.Name),
				Name:   name,
			}
			params = append(params, param)
		}

		// 递归处理嵌套的 DefineStruct
		if nestedStruct, ok := member.Type.(spec.DefineStruct); ok {
			params = append(params, extractQueryParams(nestedStruct)...)
		}
	}
	return params
}
