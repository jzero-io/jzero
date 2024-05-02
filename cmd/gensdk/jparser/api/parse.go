package api

import (
	"strings"

	"github.com/jzero-io/jzero/cmd/gensdk/vars"
	"github.com/jzero-io/jzero/daemon/pkg/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
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
	return params, nil
}

func CreateQueryParams(route *spec.Route) []*vars.QueryParam {
	if route.ResponseType == nil {
		return nil
	}

	if _, ok := route.RequestType.(spec.DefineStruct); !ok {
		return nil
	}

	members := route.RequestType.(spec.DefineStruct).GetTagMembers("form")
	params := make([]*vars.QueryParam, 0)
	for _, member := range members {
		name, _ := member.GetPropertyName()
		param := &vars.QueryParam{
			GoName: stringx.FirstUpper(member.Name),
			Name:   name,
		}
		params = append(params, param)
	}
	return params
}
