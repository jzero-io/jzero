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

func Run(args []string) error {
	baseDir := filepath.Join("desc", "proto")

	protoName := args[0]

	if strings.HasSuffix(protoName, ".proto") {
		protoName = strings.TrimSuffix(protoName, ".proto")
	}

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
		return err
	}

	if config.C.Add.Output == "file" {
		if filex.FileExists(filepath.Join(baseDir, protoName+".proto")) {
			return fmt.Errorf("%s already exists", protoName)
		}

		_ = os.MkdirAll(filepath.Dir(filepath.Join(baseDir, protoName)), 0o755)

		err = os.WriteFile(filepath.Join(baseDir, protoName+".proto"), template, 0o644)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println(string(template))
	return nil
}
