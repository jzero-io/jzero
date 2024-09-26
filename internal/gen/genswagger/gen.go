package genswagger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen"
)

func Gen(gc config.GenConfig) error {
	if pathx.FileExists(gc.Swagger.ApiDir) {
		_ = os.MkdirAll(gc.Swagger.Output, 0o755)

		if !pathx.FileExists(gc.Swagger.Output) {
			_ = os.MkdirAll(gc.Swagger.Output, 0o755)
		}

		// gen swagger by desc/api
		files, err := gen.FindRouteApiFiles(gc.Swagger.ApiDir)
		if err != nil {
			return err
		}
		for _, v := range files {
			parse, err := parser.Parse(v)
			if err != nil {
				return err
			}
			if goPackage, ok := parse.Info.Properties["go_package"]; ok {
				apiFile := fmt.Sprintf("%s.swagger.json", strings.ReplaceAll(goPackage, "/", "-"))
				cmd := exec.Command("goctl", "api", "plugin", "-plugin", "goctl-swagger=swagger -filename "+apiFile+" --schemes http", "-api", v, "-dir", gc.Swagger.Output)
				resp, err := cmd.CombinedOutput()
				if err != nil {
					return errors.Wrap(err, strings.TrimRight(string(resp), "\r\n"))
				}
				if strings.TrimRight(string(resp), "\r\n") != "" {
					fmt.Println(strings.TrimRight(string(resp), "\r\n"))
				}
			} else {
				console.Warning("[warning]: 暂不支持非 package api")
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
