package addproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func Run(args []string) (string, error) {
	baseDir := filepath.Join("desc", "proto")

	protoName := args[0]

	if strings.HasSuffix(protoName, ".proto") {
		protoName = strings.TrimSuffix(protoName, ".proto")
	}

	target := filepath.Join(baseDir, protoName+".proto")

	frameType, _ := desc.GetFrameType()
	if frameType == "" {
		frameType = "rpc"
	}

	var template []byte

	template, err := templatex.ParseTemplate(filepath.Join(frameType, "template.proto.tpl"), map[string]any{
		"Package": protoName,
		"Service": stringx.ToCamel(protoName),
	}, embeded.ReadTemplateFile(filepath.Join(frameType, "template.proto.tpl")))
	if err != nil {
		return target, err
	}

	if config.C.Add.Output == "file" {
		if filex.FileExists(target) {
			return target, fmt.Errorf("%s already exists", protoName)
		}

		_ = os.MkdirAll(filepath.Dir(filepath.Join(baseDir, protoName)), 0o755)

		err = os.WriteFile(target, template, 0o644)
		if err != nil {
			return target, err
		}
		return target, nil
	}
	fmt.Println(string(template))
	return target, nil
}
