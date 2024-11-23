package ivmaddproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type Method struct {
	Name string
	Verb string
}

func Run(c config.Config) error {
	var methods []Method
	for _, v := range c.Ivm.Add.Proto.Methods {
		split := strings.Split(v, ":")
		var method Method
		if len(split) == 2 {
			method.Name = split[1]
			method.Verb = split[0]
		} else if len(split) == 1 {
			method.Name = split[0]
			method.Verb = "get"
		} else {
			continue
		}
		methods = append(methods, method)
	}

	var version string
	var versionSuffix string

	if c.Ivm.Version == "v1" {
		version = ""
		versionSuffix = ""
	} else {
		version = c.Ivm.Version
		versionSuffix = "_" + c.Ivm.Version
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Package":    c.Ivm.Add.Proto.Name,
		"Methods":    methods,
		"Services":   c.Ivm.Add.Proto.Services,
		"Version":    version,
		"UrlVersion": c.Ivm.Version,
		"ProtoPath":  filepath.Join(c.Ivm.Version, fmt.Sprintf("%s%s.proto", c.Ivm.Add.Proto.Name, versionSuffix)),
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.proto.tpl")))
	if err != nil {
		return err
	}

	// format proto
	// Create a new printer
	// printer := &protoprint.Printer{}

	output := filepath.Join("desc", "proto", c.Ivm.Version, fmt.Sprintf("%s%s.proto", c.Ivm.Add.Proto.Name, versionSuffix))

	if pathx.FileExists(output) {
		return errors.New("proto file already exists")
	}

	return os.WriteFile(output, template, 0o644)
}
