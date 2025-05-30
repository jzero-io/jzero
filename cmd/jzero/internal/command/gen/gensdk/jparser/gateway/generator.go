package gateway

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/vars"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
)

var toCamelCaseRe = regexp.MustCompile(`(^[A-Za-z])|(_|\.)([A-Za-z])`)

func toCamelCase(str string) string {
	return toCamelCaseRe.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, "_", ""))
	})
}

func PathParam(pattern string) ([]*vars.PathParam, error) {
	if !strings.HasPrefix(pattern, "/") {
		return nil, errors.New("no leading /")
	}
	tokens, _ := tokenize(pattern[1:])

	p := parser{tokens: tokens}
	segs, err := p.topLevelSegments()
	if err != nil {
		return nil, err
	}

	params := make([]*vars.PathParam, 0)
	for i, seg := range segs {
		if v, ok := seg.(variable); ok {
			params = append(params, &vars.PathParam{
				Index:  i + 1,
				Name:   v.path,
				GoName: toCamelCase(v.path),
			})
		}
	}

	sort.Slice(params, func(i, j int) bool {
		a := params[i]
		b := params[j]
		if len(strings.Split(a.Name, ".")) < len(strings.Split(b.Name, ".")) {
			return true
		}
		return params[i].Name < params[j].Name
	})

	return params, nil
}

func CreateQueryParams(method *desc.MethodDescriptor) []*vars.QueryParam {
	queryParams := make([]*vars.QueryParam, 0)

	var f func(parent *vars.QueryParam, fields []*desc.FieldDescriptor)

	f = func(parent *vars.QueryParam, fields []*desc.FieldDescriptor) {
		for _, field := range fields {
			if field.UnwrapField().Kind() == protoreflect.MessageKind {
				q := &vars.QueryParam{
					// Field:  field,
					GoName: stringx.FirstUpper(fmt.Sprintf("%s.", field.GetName())),
					Name:   fmt.Sprintf("%s.", field.GetName()),
				}
				f(q, field.GetMessageType().GetFields())
				continue
			}
			queryParams = append(queryParams, &vars.QueryParam{
				// Field:  field,
				GoName: stringx.FirstUpper(fmt.Sprintf("%s%s", parent.GoName, field.GetName())),
				Name:   fmt.Sprintf("%s%s", parent.Name, field.GetName()),
			})
		}
	}

	f(&vars.QueryParam{GoName: "", Name: ""}, method.GetInputType().GetFields())

	return queryParams
}
