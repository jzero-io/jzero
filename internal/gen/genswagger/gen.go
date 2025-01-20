package genswagger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg/desc"
)

func Gen() error {
	if pathx.FileExists(config.C.ApiDir()) {
		_ = os.MkdirAll(config.C.Gen.Swagger.Output, 0o755)

		if !pathx.FileExists(config.C.Gen.Swagger.Output) {
			_ = os.MkdirAll(config.C.Gen.Swagger.Output, 0o755)
		}

		// gen swagger by desc/api
		files, err := desc.FindRouteApiFiles(config.C.ApiDir())
		if err != nil {
			return err
		}

		var eg errgroup.Group
		eg.SetLimit(len(files))
		for _, v := range files {
			cv := v
			eg.Go(func() error {
				parse, err := parser.Parse(cv)
				if err != nil {
					return err
				}
				apiFile := fmt.Sprintf("%s.swagger.json", strings.TrimSuffix(filepath.Base(v), filepath.Base(filepath.Ext(v))))
				if goPackage, ok := parse.Info.Properties["go_package"]; ok {
					apiFile = fmt.Sprintf("%s.swagger.json", strings.ReplaceAll(goPackage, "/", "-"))
				}
				cmd := exec.Command("goctl", "api", "plugin", "-plugin", "goctl-swagger=swagger -filename "+apiFile+" --schemes http", "-api", cv, "-dir", config.C.Gen.Swagger.Output)
				resp, err := cmd.CombinedOutput()
				if err != nil {
					return errors.Wrap(err, strings.TrimRight(string(resp), "\r\n"))
				}
				if strings.TrimRight(string(resp), "\r\n") != "" {
					fmt.Println(strings.TrimRight(string(resp), "\r\n"))
				}
				return nil
			})
		}
		if err = eg.Wait(); err != nil {
			return err
		}
	}

	if pathx.FileExists(config.C.ProtoDir()) {
		_ = os.MkdirAll(config.C.Gen.Swagger.Output, 0o755)
		protoFilepath, err := desc.GetProtoFilepath(config.C.ProtoDir())
		if err != nil {
			return err
		}

		for _, path := range protoFilepath {
			command := fmt.Sprintf("protoc -I%s -I%s %s --openapiv2_out=%s",
				config.C.ProtoDir(),
				filepath.Join(config.C.ProtoDir(), "third_party"),
				path,
				config.C.Gen.Swagger.Output,
			)
			_, err := execx.Run(command, config.C.Wd())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
