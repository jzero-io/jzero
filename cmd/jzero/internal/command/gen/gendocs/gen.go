package gendocs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
)

func Gen() (err error) {
	if pathx.FileExists(config.C.ProtoDir()) {
		_ = os.MkdirAll(config.C.Gen.Docs.Output, 0o755)
		var files []string

		switch {
		case len(config.C.Gen.Docs.Desc) > 0:
			for _, v := range config.C.Gen.Docs.Desc {
				if !osx.IsDir(v) {
					if filepath.Ext(v) == ".proto" {
						files = append(files, v)
					}
				} else {
					specifiedProtoFiles, err := desc.GetProtoFilepath(v)
					if err != nil {
						return err
					}
					files = append(files, specifiedProtoFiles...)
				}
			}
		default:
			files, err = desc.GetProtoFilepath(config.C.ProtoDir())
			if err != nil {
				return err
			}
		}

		for _, v := range config.C.Gen.Docs.DescIgnore {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == v
					})
				}
			} else {
				specifiedProtoFiles, err := desc.GetProtoFilepath(v)
				if err != nil {
					return err
				}
				for _, saf := range specifiedProtoFiles {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == saf
					})
				}
			}
		}

		command := fmt.Sprintf("protoc -I%s -I%s --doc_out=%s --doc_opt=%s,index.%s %s",
			config.C.ProtoDir(),
			filepath.Join(config.C.ProtoDir(), "third_party"),
			config.C.Gen.Docs.Output,
			getFormat(config.C.Gen.Docs.Format),
			config.C.Gen.Docs.Format,
			strings.Join(files, " "),
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
