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
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen"
)

func Gen(c config.Config) error {
	if pathx.FileExists(c.Gen.Swagger.ApiDir) {
		_ = os.MkdirAll(c.Gen.Swagger.Output, 0o755)

		if !pathx.FileExists(c.Gen.Swagger.Output) {
			_ = os.MkdirAll(c.Gen.Swagger.Output, 0o755)
		}

		// gen swagger by desc/api
		files, err := gen.FindRouteApiFiles(c.Gen.Swagger.ApiDir)
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
				if goPackage, ok := parse.Info.Properties["go_package"]; ok {
					apiFile := fmt.Sprintf("%s.swagger.json", strings.ReplaceAll(goPackage, "/", "-"))
					cmd := exec.Command("goctl", "api", "plugin", "-plugin", "goctl-swagger=swagger -filename "+apiFile+" --schemes http", "-api", cv, "-dir", c.Gen.Swagger.Output)
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
				return nil
			})
		}
		if err = eg.Wait(); err != nil {
			return err
		}
	}

	if pathx.FileExists(c.Gen.Swagger.ProtoDir) {
		_ = os.MkdirAll(c.Gen.Swagger.Output, 0o755)
		protoFilepath, err := gen.GetProtoFilepath(c.Gen.Swagger.ProtoDir)
		if err != nil {
			return err
		}

		for _, path := range protoFilepath {
			command := fmt.Sprintf("protoc -I%s -I%s %s --openapiv2_out=%s",
				c.Gen.Swagger.ProtoDir,
				filepath.Join(c.Gen.Swagger.ProtoDir, "third_party"),
				path,
				c.Gen.Swagger.Output,
			)
			_, err := execx.Run(command, c.Wd())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
