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

func AddProto(ic config.IvmConfig) error {
	var methods []Method
	for _, v := range ic.Add.Proto.Methods {
		split := strings.Split(v, ":")
		var method Method
		if len(split) == 2 {
			method.Name = split[0]
			method.Verb = split[1]
		} else if len(split) == 1 {
			method.Name = split[0]
		} else {
			continue
		}
		methods = append(methods, method)
	}

	var version string
	var versionSuffix string

	if ic.Version == "v1" {
		version = ""
		versionSuffix = ""
	} else {
		version = ic.Version
		versionSuffix = "_" + ic.Version
	}

	template, err := templatex.ParseTemplate(map[string]interface{}{
		"Package":    ic.Add.Proto.Name,
		"Methods":    methods,
		"Services":   ic.Add.Proto.Services,
		"Version":    version,
		"UrlVersion": ic.Version,
		"ProtoPath":  filepath.Join(ic.Version, fmt.Sprintf("%s%s.proto", ic.Add.Proto.Name, versionSuffix)),
	}, embeded.ReadTemplateFile(filepath.Join("ivm", "add", "template.proto.tpl")))
	if err != nil {
		return err
	}

	// format proto
	// Create a new printer
	// printer := &protoprint.Printer{}

	output := filepath.Join("desc", "proto", ic.Version, fmt.Sprintf("%s%s.proto", ic.Add.Proto.Name, versionSuffix))

	if pathx.FileExists(output) {
		return errors.New("proto file already exists")
	}

	return os.WriteFile(output, template, 0o644)
}
