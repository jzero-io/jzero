package ivmaddproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/ivm"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Methods  []string
	Name     string
	Services []string
)

type Method struct {
	Name string
	Verb string
}

func AddProto(command *cobra.Command, args []string) error {
	var methods []Method
	for _, v := range Methods {
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

	if ivm.Version == "v1" {
		version = ""
		versionSuffix = ""
	} else {
		version = ivm.Version
		versionSuffix = "_" + ivm.Version
	}

	template, err := templatex.ParseTemplate(map[string]interface{}{
		"Package":    Name,
		"Methods":    methods,
		"Services":   Services,
		"Version":    version,
		"UrlVersion": ivm.Version,
		"ProtoPath":  filepath.Join(ivm.Version, fmt.Sprintf("%s%s.proto", Name, versionSuffix)),
	}, embeded.ReadTemplateFile("template.proto.tpl"))
	if err != nil {
		return err
	}

	// format proto
	// Create a new printer
	// printer := &protoprint.Printer{}

	output := filepath.Join("desc", "proto", ivm.Version, fmt.Sprintf("%s%s.proto", Name, versionSuffix))

	if pathx.FileExists(output) {
		return errors.New("proto file already exists")
	}

	return os.WriteFile(output, template, 0o644)
}
