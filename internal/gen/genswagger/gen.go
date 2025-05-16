package genswagger

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/jaronnie/genius"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/osx"
	"github.com/jzero-io/jzero/pkg/stringx"
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

			if config.C.GoctlVersion().GreaterThanOrEqual(version.Must(version.NewVersion("1.8.3"))) {
				eg.Go(func() error {
					parse, err := apiparser.Parse(cv, nil)
					if err != nil {
						return err
					}

					relativePath := strings.TrimPrefix(cv, config.C.ApiDir())
					relativePath = strings.TrimPrefix(relativePath, "/")
					pathBasedName := strings.ReplaceAll(relativePath, "/", "-")
					pathBasedName = strings.TrimSuffix(pathBasedName, ".api")

					apiFile := fmt.Sprintf("%s.swagger", pathBasedName)
					goPackage, ok := parse.Info.Properties["go_package"]
					if ok && goPackage != "" {
						apiFile = fmt.Sprintf("%s.swagger", strings.ReplaceAll(goPackage, "/", "-"))
					}

					if apiFile == "version.swagger" {
						apiFile = "swagger"
					}

					cmd := exec.Command("goctl", "api", "swagger", "--api", cv, "--filename", apiFile, "--dir", config.C.Gen.Swagger.Output)

					logx.Debug(cmd.String())
					resp, err := cmd.CombinedOutput()
					if err != nil {
						return errors.Wrap(err, strings.TrimRight(string(resp), "\r\n"))
					}
					if strings.TrimRight(string(resp), "\r\n") != "" {
						fmt.Println(strings.TrimRight(string(resp), "\r\n"))
					}

					// 兼容处理
					file, err := os.ReadFile(filepath.Join(config.C.Gen.Swagger.Output, apiFile+".json"))
					if err != nil {
						return err
					}
					g, err := genius.NewFromRawJSON(file)
					if err != nil {
						return err
					}

					// 处理 host 值
					if cast.ToString(g.Get("host")) == "127.0.0.1" {
						_ = g.Set("host", "")
					}

					// 处理 securityDefinitions 值
					if g.Get("securityDefinitions") == nil {
						_ = g.Set("securityDefinitions", map[string]any{
							"apiKey": map[string]any{
								"type":        "apiKey",
								"description": "Enter Authorization",
								"name":        "Authorization",
								"in":          "header",
							},
						})
					}

					// 处理 schemes 值
					if len(cast.ToStringSlice(g.Get("schemes"))) == 1 && cast.ToStringSlice(g.Get("schemes"))[0] == "https" {
						_ = g.Set("schemes", []string{"http", "https"})
					}

					pathMaps := cast.ToStringMap(g.Get("paths"))
					for pmk := range pathMaps {
						pathMethodsMap := cast.ToStringMap(pathMaps[pmk])
						for pmmk := range pathMethodsMap {
							for _, group := range parse.Service.Groups {
								for _, route := range group.Routes {
									logx.Debugf("get route prefix: %s", route.GetAnnotation("prefix"))
									if group.GetAnnotation("prefix") != "" {
										route.Path = group.GetAnnotation("prefix") + route.Path
									}
									if route.Method == pmmk && route.Path == adjustHttpPath(pmk) && group.GetAnnotation("group") != "" {
										h := strings.TrimSuffix(route.Handler, "Handler")
										groupName := group.GetAnnotation("group")

										if config.C.Gen.Swagger.Route2Code || config.C.Gen.Route2Code {
											_ = g.Set(fmt.Sprintf("paths.%s.%s.description", pmk, pmmk), "接口权限编码"+":"+stringx.FirstLower(strings.ReplaceAll(groupName, "/", ":"))+":"+stringx.FirstLower(h))
										}
									}
								}
							}

							// 处理 tags
							tags := cast.ToStringSlice(g.Get(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk)))
							if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
								if goPackage != "" {
									_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{goPackage})
								} else {
									for _, group := range parse.Service.Groups {
										for _, route := range group.Routes {
											logx.Debugf("get route prefix: %s", route.GetAnnotation("prefix"))
											if group.GetAnnotation("prefix") != "" {
												route.Path = group.GetAnnotation("prefix") + route.Path
											}
											if route.Method == pmmk && route.Path == adjustHttpPath(pmk) && group.GetAnnotation("group") != "" {
												_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{group.GetAnnotation("group")})
												break
											}
										}
									}
								}
							}

							// 处理 security
							/*
								"security": [
								          {
								            "apiKey": []
								          }
								        ],
							*/
							if g.Get(fmt.Sprintf("paths.%s.%s.security", pmk, pmmk)) == nil {
								_ = g.Set(fmt.Sprintf("paths.%s.%s.security", pmk, pmmk), []map[string][]any{
									{
										"apiKey": []any{},
									},
								})
							}
						}
					}

					encodeToJSON, err := g.EncodeToPrettyJSON()
					if err != nil {
						return err
					}
					err = os.WriteFile(filepath.Join(config.C.Gen.Swagger.Output, apiFile+".json"), encodeToJSON, 0o644)
					if err != nil {
						return err
					}
					return nil
				})

				if err = eg.Wait(); err != nil {
					return err
				}
			} else {
				eg.Go(func() error {
					parse, err := apiparser.Parse(cv, nil)
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

				if err = eg.Wait(); err != nil {
					return err
				}
			}
		}

		// merge swagger to one file swagger.json
		// use swagger.api to set global config
		if config.C.Gen.Swagger.Merge && config.C.GoctlVersion().GreaterThanOrEqual(version.Must(version.NewVersion("1.8.3"))) {
			swaggerJson, err := os.ReadFile(filepath.Join(config.C.SwaggerDir(), "swagger.json"))
			if err != nil {
				return err
			}
			swaggerJsonG, err := genius.NewFromRawJSON(swaggerJson)
			if err != nil {
				return err
			}

			dir, err := os.ReadDir(config.C.SwaggerDir())
			if err == nil {
				for _, sj := range dir {
					if sj.Name() != "swagger.json" {
						file, err := os.ReadFile(filepath.Join(config.C.SwaggerDir(), sj.Name()))
						if err == nil {
							g, err := genius.NewFromRawJSON(file)
							if err == nil {
								paths := g.Get("paths")
								pathsMarshal, _ := json.Marshal(paths)

								pathMaps := make(map[string]any)
								_ = json.Unmarshal(pathsMarshal, &pathMaps)

								for pmk, pmv := range pathMaps {
									_ = swaggerJsonG.Set(fmt.Sprintf("paths.%s", pmk), pmv)
								}
							}
						}
					}
				}

				encodeToJSON, err := swaggerJsonG.EncodeToPrettyJSON()
				if err != nil {
					return err
				}
				err = os.WriteFile(filepath.Join(config.C.SwaggerDir(), "swagger.json"), encodeToJSON, 0o644)
				if err != nil {
					return err
				}
			}
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

func adjustHttpPath(path string) string {
	path = strings.ReplaceAll(path, "{", ":")
	path = strings.ReplaceAll(path, "}", "")
	return path
}
