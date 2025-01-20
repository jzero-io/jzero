package genrpc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero-contrib/filex"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/config"
)

func (jr *JzeroRpc) patchSvc() error {
	namingFormat, err := format.FileNamingFormat(config.C.Gen.Style, "service_context.go")
	if err != nil {
		return err
	}

	if !filex.DirExists(filepath.Join("internal", "svc")) || filex.FileExists(filepath.Join("internal", "svc", namingFormat)) {
		return nil
	}

	dir, err := os.ReadDir(filepath.Join("internal", "svc"))
	if err != nil {
		return err
	}
	for _, v := range dir {
		if !v.IsDir() {
			if strings.HasPrefix(v.Name(), "service") && strings.HasSuffix(v.Name(), "context.go") {
				if err = os.Rename(filepath.Join("internal", "svc", v.Name()), filepath.Join("internal", "svc", namingFormat)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
