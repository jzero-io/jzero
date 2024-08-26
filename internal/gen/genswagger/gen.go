package genswagger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func Gen(gc config.GenConfig) error {
	if pathx.FileExists(gc.Swagger.ApiDir) {
		_ = os.MkdirAll(gc.Swagger.Output, 0o755)
		mainApiFile, isDelete, err := gen.GetMainApiFilePath(gc.Swagger.ApiDir)
		if err != nil {
			return err
		}
		defer func() {
			if isDelete {
				_ = os.Remove(mainApiFile)
			}
		}()

		if !pathx.FileExists(gc.Swagger.Output) {
			_ = os.MkdirAll(gc.Swagger.Output, 0o755)
		}

		// gen swagger by desc/api
		if mainApiFile != "" {
			apiFile := fmt.Sprintf("%s.swagger.json", gen.GetApiServiceName(gc.Swagger.ApiDir))
			cmd := exec.Command("goctl", "api", "plugin", "-plugin", "goctl-swagger=swagger -filename "+apiFile+" --schemes http", "-api", mainApiFile, "-dir", gc.Swagger.Output)
			resp, err := cmd.CombinedOutput()
			if err != nil {
				return errors.Wrap(err, strings.TrimRight(string(resp), "\r\n"))
			}
			if strings.TrimRight(string(resp), "\r\n") != "" {
				fmt.Println(strings.TrimRight(string(resp), "\r\n"))
			}
		}
	}

	if pathx.FileExists(gc.Swagger.ProtoDir) {
		_ = os.MkdirAll(gc.Swagger.Output, 0o755)
		protoFilepath, err := gen.GetProtoFilepath(gc.Swagger.ProtoDir)
		if err != nil {
			return err
		}

		for _, path := range protoFilepath {
			command := fmt.Sprintf("protoc -I%s -I%s %s --openapiv2_out=%s",
				gc.Swagger.ProtoDir,
				filepath.Join(gc.Swagger.ProtoDir, "third_party"),
				path,
				gc.Swagger.Output,
			)
			_, err := execx.Run(command, gc.Swagger.Wd())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
