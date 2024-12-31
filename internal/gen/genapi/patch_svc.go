package genapi

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero-contrib/filex"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

func (ja *JzeroApi) patchSvc() error {
	namingFormat, err := format.FileNamingFormat(ja.Style, "service_context.go")
	if err != nil {
		return err
	}
	if filex.FileExists(filepath.Join("internal", "svc")) {
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
