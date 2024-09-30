package gendocs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg/desc"
)

func Gen(c config.Config) error {
	if pathx.FileExists(c.Gen.Docs.ProtoDir) {
		_ = os.MkdirAll(c.Gen.Docs.Output, 0o755)
		protoFilepath, err := desc.GetProtoFilepath(c.Gen.Swagger.ProtoDir)
		if err != nil {
			return err
		}

		command := fmt.Sprintf("protoc -I%s -I%s --doc_out=%s --doc_opt=%s,index.%s %s",
			c.Gen.Docs.ProtoDir,
			filepath.Join(c.Gen.Docs.ProtoDir, "third_party"),
			c.Gen.Docs.Output,
			getFormat(c.Gen.Docs.Format),
			c.Gen.Docs.Format,
			strings.Join(protoFilepath, " "),
		)
		_, err = execx.Run(command, c.Wd())
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
