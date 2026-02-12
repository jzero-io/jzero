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

				// 保持目录结构，生成到 desc/swagger 下
				var relPath string

				// 检查是否是插件文件
				pluginName := getPluginNameFromFilePath(v)
				if pluginName != "" {
					// 插件文件处理：找到 desc/api 在路径中的位置
					descApiPath := filepath.Join("desc", "api") + string(filepath.Separator)
					descApiIndex := strings.Index(v, descApiPath)
					var pluginApiDir string
					if descApiIndex == -1 {
						// 如果找不到 desc/api 模式，尝试查找路径末尾是否以 desc/api 结尾
						if strings.HasSuffix(filepath.Dir(v), filepath.Join("desc", "api")) {
							pluginApiDir = filepath.Dir(v)
						} else {
							return fmt.Errorf("invalid plugin api path: %s", v)
						}
					} else {
						pluginApiDir = v[:descApiIndex+len(descApiPath)]
					}

					var relErr error
					relPath, relErr = filepath.Rel(pluginApiDir, v)
					if relErr != nil {
						return relErr
					}
					// 在插件目录下保持结构
					relPath = filepath.Join("plugins", pluginName, relPath)
				} else {
					// 普通 API 文件处理
					relPath, err = filepath.Rel(config.C.ApiDir(), v)
					if err != nil {
						return err
					}
				}

				// 将 .api 扩展名替换为 .swagger
				swaggerFileName := strings.TrimSuffix(relPath, ".api") + ".swagger"

				// 创建输出目录结构
				outputDir := filepath.Join(config.C.Gen.Swagger.Output, filepath.Dir(swaggerFileName))
				if err := os.MkdirAll(outputDir, 0o755); err != nil {
					return err
				}

				apiFile := filepath.Base(swaggerFileName)
				goPackage := parse.Info.Properties["go_package"]

				cmd := exec.Command("goctl", "api", "swagger", "--api", v, "--filename", apiFile, "--dir", outputDir)

				logx.Debug(cmd.String())
				resp, err := cmd.CombinedOutput()
				if err != nil {
					return errors.Wrap(err, strings.TrimRight(string(resp), "\r\n"))
				}
				if strings.TrimRight(string(resp), "\r\n") != "" {
					fmt.Println(strings.TrimRight(string(resp), "\r\n"))
				}

				// 兼容处理
				file, err := os.ReadFile(filepath.Join(outputDir, apiFile+".json"))
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
						pluginName = getPluginNameFromFilePath(v)

						if pluginName != "" {
							// 插件文件：处理已存在的 tags，为每个 tag 添加插件前缀
							if len(tags) > 0 && !(len(tags) == 1 && tags[0] == "") {
								var newTags []string
								for _, tag := range tags {
									if tag != "" {
										newTags = append(newTags, "plugins/"+pluginName+"/"+tag)
									}
								}
								if len(newTags) > 0 {
									_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), newTags)
								}
							} else {
								// 如果没有 tags，设置默认 tags
								if goPackage != "" {
									tagValue := "plugins/" + pluginName + "/" + goPackage
									_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{tagValue})
								} else {
									for _, group := range parse.Service.Groups {
										for _, route := range group.Routes {
											logx.Debugf("get route prefix: %s", route.GetAnnotation("prefix"))
											if group.GetAnnotation("prefix") != "" {
												route.Path = group.GetAnnotation("prefix") + route.Path
											}
											if route.Method == pmmk && route.Path == adjustHttpPath(pmk) && group.GetAnnotation("group") != "" {
												tagValue := "plugins/" + pluginName + "/" + group.GetAnnotation("group")
												_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{tagValue})
												break
											}
										}
									}
								}
							}
						} else {
							// 普通文件：只在没有 tags 时设置默认值
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
						}

						// 处理 operationId
						pluginName = getPluginNameFromFilePath(v)
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
				err = os.WriteFile(filepath.Join(outputDir, apiFile+".json"), encodeToJSON, 0o644)
				if err != nil {
					return err
				}
				return nil
			})

			if err = eg.Wait(); err != nil {
				return err
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
					specifiedProtoFiles, err := desc.FindRpcServiceProtoFiles(v)
					if err != nil {
						return err
					}
					files = append(files, specifiedProtoFiles...)
				}
			}
		default:
			files, err = desc.FindRpcServiceProtoFiles(config.C.ProtoDir())
			if err != nil {
				return err
			}

			// 增加 plugins 的 proto 文件
			plugins, err := plugin.GetPlugins()
			if err == nil {
				for _, p := range plugins {
					if pathx.FileExists(filepath.Join(p.Path, "desc", "proto")) {
						pluginFiles, err := desc.FindRpcServiceProtoFiles(filepath.Join(p.Path, "desc", "proto"))
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
				if filepath.Ext(v) == ".proto" {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == v
					})
				}
			} else {
				specifiedProtoFiles, err := desc.FindRpcServiceProtoFiles(v)
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
			// 检查是否是插件文件
			pluginName := getPluginNameFromFilePath(path)
			var pluginProtoDir string
			var outputDir string

			if pluginName != "" {
				// 插件文件处理：找到 proto 目录在路径中的位置
				protoPath := filepath.Join("", config.C.ProtoDir()) + string(filepath.Separator)
				descProtoIndex := strings.Index(path, protoPath)
				if descProtoIndex == -1 {
					// 如果找不到 proto 路径模式，尝试查找路径末尾是否以 proto 目录结尾
					if strings.HasSuffix(filepath.Dir(path), filepath.Join("", config.C.ProtoDir())) {
						pluginProtoDir = filepath.Dir(path)
					} else {
						return fmt.Errorf("invalid plugin proto path: %s", path)
					}
				} else {
					pluginProtoDir = path[:descProtoIndex+len(protoPath)]
				}
				// 插件文件直接生成到 plugins/{插件名}/ 目录下
				outputDir = filepath.Join(config.C.Gen.Swagger.Output, "plugins", pluginName)
			} else {
				// 普通 Proto 文件直接生成到根目录下
				outputDir = config.C.Gen.Swagger.Output
			}

			// 创建输出目录结构
			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return err
			}

			// 为插件文件添加插件路径到 protoc 的 -I 参数
			var includeArgs []string
			if pluginName != "" {
				includeArgs = append(includeArgs, "-I"+pluginProtoDir)
				// 还需要包含插件的 third_party 目录
				pluginThirdParty := filepath.Join(pluginProtoDir, "third_party")
				if pathx.FileExists(pluginThirdParty) {
					includeArgs = append(includeArgs, "-I"+pluginThirdParty)
				}
			}
			includeArgs = append(includeArgs, "-I"+config.C.ProtoDir())
			includeArgs = append(includeArgs, "-I"+filepath.Join(config.C.ProtoDir(), "third_party"))

			command := fmt.Sprintf("protoc %s %s --openapiv2_out=%s",
				strings.Join(includeArgs, " "),
				path,
				outputDir,
			)
			_, err = execx.Run(command, config.C.Wd())
			if err != nil {
				return err
			}
		}
	}

	// 统一的 merge 处理，合并所有生成的 swagger 文件
	if config.C.Gen.Swagger.Merge {
		err = mergeSwaggerFiles()
		if err != nil {
			return err
		}
	}

	return nil
}

// mergeSwaggerFiles 递归扫描并合并所有的 swagger 文件
func mergeSwaggerFiles() error {
	swaggerJson := embeded.ReadTemplateFile(filepath.Join("swagger", "swagger.json.tpl"))

	swaggerJsonG, err := genius.NewFromRawJSON(swaggerJson)
	if err != nil {
		return err
	}

	// 递归扫描所有 swagger 文件
	swaggerFiles, err := findAllSwaggerFiles(config.C.Gen.Swagger.Output)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	// 合并所有文件的 paths
	for _, filePath := range swaggerFiles {
		// 跳过主 swagger.json 文件
		if filepath.Base(filePath) == "swagger.json" {
			continue
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			logx.Errorf("failed to read swagger file %s: %v", filePath, err)
			continue
		}

		g, err := genius.NewFromRawJSON(file)
		if err != nil {
			logx.Errorf("failed to parse swagger file %s: %v", filePath, err)
			continue
		}

		// 合并 paths
		paths := g.Get("paths")
		if paths != nil {
			pathsMarshal, _ := json.Marshal(paths)
			pathMaps := make(map[string]any)
			_ = json.Unmarshal(pathsMarshal, &pathMaps)

			for pmk, pmv := range pathMaps {
				_ = swaggerJsonG.Set(fmt.Sprintf("paths.%s", pmk), pmv)
			}
		}

		// 合并 definitions（如果存在）
		definitions := g.Get("definitions")
		if definitions != nil {
			definitionsMarshal, _ := json.Marshal(definitions)
			definitionsMap := make(map[string]any)
			_ = json.Unmarshal(definitionsMarshal, &definitionsMap)

			for defKey, defValue := range definitionsMap {
				_ = swaggerJsonG.Set(fmt.Sprintf("definitions.%s", defKey), defValue)
			}
		}
	}

	// 写入合并后的文件
	encodeToJSON, err := swaggerJsonG.EncodeToPrettyJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(config.C.Gen.Swagger.Output, "swagger.json"), encodeToJSON, 0o644)
}

// findAllSwaggerFiles 递归查找所有的 swagger JSON 文件
func findAllSwaggerFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .json 文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
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
