package ivmaddproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type Method struct {
	Name string
	Verb string
}

func Run() error {
	var methods []Method
	for _, v := range config.C.Ivm.Add.Proto.Methods {
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

	if config.C.Ivm.Version == "v1" {
		version = ""
		versionSuffix = ""
	} else {
		version = config.C.Ivm.Version
		versionSuffix = "_" + config.C.Ivm.Version
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Package":    config.C.Ivm.Add.Proto.Name,
		"Methods":    methods,
		"Services":   config.C.Ivm.Add.Proto.Services,
		"Version":    version,
		"UrlVersion": config.C.Ivm.Version,
		"ProtoPath":  filepath.Join(config.C.Ivm.Version, fmt.Sprintf("%s%s.proto", config.C.Ivm.Add.Proto.Name, versionSuffix)),
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.proto.tpl")))
	if err != nil {
		return err
	}

	// format proto
	// Create a new printer
	// printer := &protoprint.Printer{}

	output := filepath.Join("desc", "proto", config.C.Ivm.Version, fmt.Sprintf("%s%s.proto", config.C.Ivm.Add.Proto.Name, versionSuffix))

	if pathx.FileExists(output) {
		return errors.New("proto file already exists")
	}

	return os.WriteFile(output, template, 0o644)
}
