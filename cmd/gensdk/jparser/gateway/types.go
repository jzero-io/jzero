// Refer: https://github.com/grpc-ecosystem/grpc-gateway

package gateway

import (
	"fmt"
	"strings"
)

type segment interface {
	fmt.Stringer
}

type wildcard struct{}

type deepWildcard struct{}

type literal string

type variable struct {
	path     string
	segments []segment
}

func (wildcard) String() string {
	return "*"
}

func (deepWildcard) String() string {
	return "**"
}

func (l literal) String() string {
	return string(l)
}

func (v variable) String() string {
	var segs []string
	for _, s := range v.segments {
		segs = append(segs, s.String())
	}
	return fmt.Sprintf("{%s=%s}", v.path, strings.Join(segs, "/"))
}
