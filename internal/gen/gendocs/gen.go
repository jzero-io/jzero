package gendocs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen"
)

func Gen(gc config.GenConfig) error {
	if pathx.FileExists(gc.Docs.ProtoDir) {
		_ = os.MkdirAll(gc.Docs.Output, 0o755)
		protoFilepath, err := gen.GetProtoFilepath(gc.Swagger.ProtoDir)
		if err != nil {
			return err
		}

		command := fmt.Sprintf("protoc -I%s -I%s --doc_out=%s --doc_opt=%s,index.%s %s",
			gc.Docs.ProtoDir,
			filepath.Join(gc.Docs.ProtoDir, "third_party"),
			gc.Docs.Output,
			getFormat(gc.Docs.Format),
			gc.Docs.Format,
			strings.Join(protoFilepath, " "),
		)
		_, err = execx.Run(command, gc.Swagger.Wd())
		if err != nil {
			return err
		}
	}

	return nil
}

func getFormat(format string) string {
	switch format {
	case "md":
		return "markdown"
	}
	return format
}
