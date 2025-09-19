package genswagger

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jaronnie/genius"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/plugin"
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

			// 增加 plugins 的 api 文件
			plugins, err := plugin.GetPlugins()
			if err == nil {
				for _, p := range plugins {
					if pathx.FileExists(filepath.Join(p.Path, "desc", "api")) {
						pluginFiles, err := desc.FindRouteApiFiles(filepath.Join(p.Path, "desc", "api"))
						if err != nil {
							return err
						}
						files = append(files, pluginFiles...)
					}
				}
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
			eg.Go(func() error {
				parse, err := apiparser.Parse(v, nil)
				if err != nil {
					return err
				}

				pathBasedName := strings.ReplaceAll(v, string(filepath.Separator), "-")
				pathBasedName = strings.ReplaceAll(pathBasedName, "desc-api-", "")
				pathBasedName = strings.TrimSuffix(pathBasedName, ".api")
				apiFile := fmt.Sprintf("%s.swagger", pathBasedName)
				goPackage := parse.Info.Properties["go_package"]

				cmd := exec.Command("goctl", "api", "swagger", "--api", v, "--filename", apiFile, "--dir", config.C.Gen.Swagger.Output)

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

				/*
				 "x-date": "",
				  "x-description": "This is a goctl generated swagger file.",
				  "x-github": "https://github.com/zeromicro/go-zero",
				  "x-go-zero-doc": "https://go-zero.dev/",
				  "x-goctl-version": "1.9.0"
				*/

				// 删除 x-date
				g.Del("x-date")
				g.Del("x-description")
				g.Del("x-github")
				g.Del("x-go-zero-doc")
				g.Del("x-goctl-version")

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
							pluginName := getPluginNameFromFilePath(v)
							if goPackage != "" {
								tagValue := goPackage
								if pluginName != "" {
									tagValue = "plugins/" + pluginName + "/" + goPackage
								}
								_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{tagValue})
							} else {
								for _, group := range parse.Service.Groups {
									for _, route := range group.Routes {
										logx.Debugf("get route prefix: %s", route.GetAnnotation("prefix"))
										if group.GetAnnotation("prefix") != "" {
											route.Path = group.GetAnnotation("prefix") + route.Path
										}
										if route.Method == pmmk && route.Path == adjustHttpPath(pmk) && group.GetAnnotation("group") != "" {
											tagValue := group.GetAnnotation("group")
											if pluginName != "" {
												tagValue = "plugins/" + pluginName + "/" + group.GetAnnotation("group")
											}
											_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{tagValue})
											break
										}
									}
								}
							}
						}

						// 处理 operationId
						pluginName := getPluginNameFromFilePath(v)
						if pluginName != "" {
							operationId := cast.ToString(g.Get(fmt.Sprintf("paths.%s.%s.operationId", pmk, pmmk)))
							if operationId != "" {
								newOperationId := "plugins/" + pluginName + "/" + operationId
								_ = g.Set(fmt.Sprintf("paths.%s.%s.operationId", pmk, pmmk), newOperationId)
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
		}

		// merge swagger to one file swagger.json
		// use swagger.api to set global config
		if config.C.Gen.Swagger.Merge {
			swaggerJson, err := os.ReadFile(filepath.Join(config.C.SwaggerDir(), "swagger.json"))
			if err != nil {
				swaggerJson = embeded.ReadTemplateFile(filepath.Join("swagger", "swagger.json.tpl"))
				if swaggerJson == nil {
					return err
				}
				err = nil
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

func getPluginNameFromFilePath(filePath string) string {
	if strings.Contains(filePath, "plugins"+string(filepath.Separator)) {
		parts := strings.Split(filePath, string(filepath.Separator))
		for i, part := range parts {
			if part == "plugins" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	return ""
}
