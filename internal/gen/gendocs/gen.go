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

func Gen() error {
	if pathx.FileExists(config.C.ProtoDir()) {
		_ = os.MkdirAll(config.C.Gen.Docs.Output, 0o755)
		protoFilepath, err := desc.GetProtoFilepath(config.C.ProtoDir())
		if err != nil {
			return err
		}

		command := fmt.Sprintf("protoc -I%s -I%s --doc_out=%s --doc_opt=%s,index.%s %s",
			config.C.ProtoDir(),
			filepath.Join(config.C.ProtoDir(), "third_party"),
			config.C.Gen.Docs.Output,
			getFormat(config.C.Gen.Docs.Format),
			config.C.Gen.Docs.Format,
			strings.Join(protoFilepath, " "),
		)
		_, err = execx.Run(command, config.C.Wd())
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
