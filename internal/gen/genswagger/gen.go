package genswagger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/osx"
)

func Gen() (err error) {
	if pathx.FileExists(config.C.ApiDir()) {
		_ = os.MkdirAll(config.C.Gen.Swagger.Output, 0o755)

		if !pathx.FileExists(config.C.Gen.Swagger.Output) {
			_ = os.MkdirAll(config.C.Gen.Swagger.Output, 0o755)
		}

		var files []string

		switch {
		case len(config.C.Gen.Swagger.Desc) > 0:
			for _, v := range config.C.Gen.Swagger.Desc {
				if !osx.IsDir(v) {
					if filepath.Ext(v) == ".api" {
						files = append(files, v)
					}
				} else {
					specifiedApiFiles, err := desc.FindApiFiles(v)
					if err != nil {
						return err
					}
					files = append(files, specifiedApiFiles...)
				}
			}
		default:
			files, err = desc.FindRouteApiFiles(config.C.ApiDir())
			if err != nil {
				return err
			}
		}

		for _, v := range config.C.Gen.Swagger.DescIgnore {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".api" {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == v
					})
				}
			} else {
				specifiedApiFiles, err := desc.FindApiFiles(v)
				if err != nil {
					return err
				}
				for _, saf := range specifiedApiFiles {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == saf
					})
				}
			}
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
				cmd := exec.Command("goctl", "api", "plugin", "-plugin", "jzero-swagger=swagger -filename "+apiFile+" --schemes http,https", "-api", cv, "-dir", config.C.Gen.Swagger.Output)
				if config.C.Gen.Route2Code || config.C.Gen.Swagger.Route2Code {
					cmd = exec.Command("goctl", "api", "plugin", "-plugin", "jzero-swagger=swagger -filename "+apiFile+" --schemes http,https "+" --route2code ", "-api", cv, "-dir", config.C.Gen.Swagger.Output)
				}

				logx.Debug(cmd.String())
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

		var files []string

		switch {
		case len(config.C.Gen.Swagger.Desc) > 0:
			for _, v := range config.C.Gen.Swagger.Desc {
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

		for _, v := range config.C.Gen.Swagger.DescIgnore {
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

		for _, path := range files {
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
