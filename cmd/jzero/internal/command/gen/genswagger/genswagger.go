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
	apiparser "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console/progress"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/stringx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/serverless"
)

func Gen() (err error) {
	showProgress := !config.C.Quiet
	hasSwaggerInput := hasSwaggerSourceInput()

	if err = executeStage(
		console.Green("Gen")+" "+console.Yellow("swagger api"),
		false,
		showProgress,
		runRegularAPISwagger,
	); err != nil {
		return err
	}

	if err = executeStage(
		console.Green("Gen")+" "+console.Yellow("swagger plugin api"),
		false,
		showProgress,
		runPluginAPISwagger,
	); err != nil {
		return err
	}

	if err = executeStage(
		console.Green("Gen")+" "+console.Yellow("swagger proto"),
		false,
		showProgress,
		runRegularProtoSwagger,
	); err != nil {
		return err
	}

	if err = executeStage(
		console.Green("Gen")+" "+console.Yellow("swagger plugin proto"),
		false,
		showProgress,
		runPluginProtoSwagger,
	); err != nil {
		return err
	}

	if config.C.Gen.Swagger.Merge && hasSwaggerInput {
		if err = executeStage(
			console.Green("Merge")+" "+console.Yellow("swagger"),
			showProgress && hasSwaggerInput,
			showProgress,
			runMergeSwagger,
		); err != nil {
			return err
		}
	}

	return nil
}

func runRegularAPISwagger(progressChan chan<- progress.Message) error {
	files, err := listSwaggerAPIFiles()
	if err != nil {
		return err
	}
	regularFiles, _ := splitPluginPaths(files)
	return runAPISwagger(regularFiles, progressChan)
}

func runPluginAPISwagger(progressChan chan<- progress.Message) error {
	files, err := listSwaggerAPIFiles()
	if err != nil {
		return err
	}
	_, pluginFiles := splitPluginPaths(files)
	return runAPISwagger(pluginFiles, progressChan)
}

func runRegularProtoSwagger(progressChan chan<- progress.Message) error {
	files, err := listSwaggerProtoFiles()
	if err != nil {
		return err
	}
	regularFiles, _ := splitPluginPaths(files)
	return runProtoSwagger(regularFiles, progressChan)
}

func runPluginProtoSwagger(progressChan chan<- progress.Message) error {
	files, err := listSwaggerProtoFiles()
	if err != nil {
		return err
	}
	_, pluginFiles := splitPluginPaths(files)
	return runProtoSwagger(pluginFiles, progressChan)
}

func executeStage(title string, headerShown, showProgress bool, fn func(chan<- progress.Message) error) error {
	if !showProgress {
		return fn(nil)
	}

	progressChan := make(chan progress.Message, 10)
	done := make(chan struct{})
	var stageErr error

	if headerShown {
		fmt.Printf("%s\n", console.BoxHeader("", title))
	}

	go func() {
		stageErr = fn(progressChan)
		close(done)
	}()

	state := progress.ConsumeStage(progressChan, done, title, false, headerShown)
	progress.FinishStage(title, false, &state, stageErr)

	if stageErr != nil {
		return console.MarkRenderedError(stageErr)
	}

	return nil
}

func runAPISwagger(files []string, progressChan chan<- progress.Message) error {
	if err := ensureSwaggerOutputDir(); err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	var eg errgroup.Group
	eg.SetLimit(len(files))
	for _, file := range files {
		file := file
		eg.Go(func() error {
			if err := processSwaggerAPIFile(file); err != nil {
				return errors.Wrapf(err, "swagger api file: %s", file)
			}
			if progressChan != nil {
				progressChan <- progress.NewFile(file)
			}
			return nil
		})
	}

	return eg.Wait()
}

func listSwaggerAPIFiles() ([]string, error) {
	var files []string

	switch {
	case len(config.C.Gen.Swagger.Desc) > 0:
		for _, v := range config.C.Gen.Swagger.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".api" {
					files = append(files, v)
				}
				continue
			}

			specifiedApiFiles, err := desc.FindApiFiles(v)
			if err != nil {
				return nil, err
			}
			files = append(files, specifiedApiFiles...)
		}
	default:
		if pathx.FileExists(config.C.ApiDir()) {
			var err error
			files, err = desc.FindRouteApiFiles(config.C.ApiDir())
			if err != nil {
				return nil, err
			}
		}

		plugins, err := serverless.GetPlugins()
		if err == nil {
			for _, p := range plugins {
				if pathx.FileExists(filepath.Join(p.Path, "desc", "api")) {
					pluginFiles, err := desc.FindRouteApiFiles(filepath.Join(p.Path, "desc", "api"))
					if err != nil {
						return nil, err
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
			continue
		}

		specifiedApiFiles, err := desc.FindApiFiles(v)
		if err != nil {
			return nil, err
		}
		for _, saf := range specifiedApiFiles {
			files = lo.Reject(files, func(item string, _ int) bool {
				return item == saf
			})
		}
	}

	return files, nil
}

func processSwaggerAPIFile(apiPath string) error {
	parse, err := apiparser.Parse(apiPath, nil)
	if err != nil {
		return err
	}

	var relPath string

	pluginName := getPluginNameFromFilePath(apiPath)
	if pluginName != "" {
		descApiPath := filepath.Join("desc", "api") + string(filepath.Separator)
		descApiIndex := strings.Index(apiPath, descApiPath)
		var pluginAPIDir string
		if descApiIndex == -1 {
			if strings.HasSuffix(filepath.Dir(apiPath), filepath.Join("desc", "api")) {
				pluginAPIDir = filepath.Dir(apiPath)
			} else {
				return fmt.Errorf("invalid plugin api path: %s", apiPath)
			}
		} else {
			pluginAPIDir = apiPath[:descApiIndex+len(descApiPath)]
		}

		relPath, err = filepath.Rel(pluginAPIDir, apiPath)
		if err != nil {
			return err
		}
		relPath = filepath.Join("plugins", pluginName, relPath)
	} else {
		relPath, err = filepath.Rel(config.C.ApiDir(), apiPath)
		if err != nil {
			return err
		}
	}

	swaggerFileName := strings.TrimSuffix(relPath, ".api") + ".swagger"
	outputDir := filepath.Join(config.C.Gen.Swagger.Output, filepath.Dir(swaggerFileName))
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	apiFile := filepath.Base(swaggerFileName)
	goPackage := parse.Info.Properties["go_package"]

	cmd := exec.Command("goctl", "api", "swagger", "--api", apiPath, "--filename", apiFile, "--dir", outputDir)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, strings.TrimRight(string(resp), "\r\n"))
	}

	file, err := os.ReadFile(filepath.Join(outputDir, apiFile+".json"))
	if err != nil {
		return err
	}
	g, err := genius.NewFromRawJSON(file)
	if err != nil {
		return err
	}

	if cast.ToString(g.Get("host")) == "127.0.0.1" {
		_ = g.Set("host", "")
	}

	g.Del("x-date")
	g.Del("x-description")
	g.Del("x-github")
	g.Del("x-go-zero-doc")
	g.Del("x-goctl-version")

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

	if len(cast.ToStringSlice(g.Get("schemes"))) == 1 && cast.ToStringSlice(g.Get("schemes"))[0] == "https" {
		_ = g.Set("schemes", []string{"http", "https"})
	}

	pathMaps := cast.ToStringMap(g.Get("paths"))
	for pmk := range pathMaps {
		pathMethodsMap := cast.ToStringMap(pathMaps[pmk])
		for pmmk := range pathMethodsMap {
			for _, group := range parse.Service.Groups {
				for _, route := range group.Routes {
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

			tags := cast.ToStringSlice(g.Get(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk)))
			pluginName = getPluginNameFromFilePath(apiPath)

			if pluginName != "" {
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
					if goPackage != "" {
						tagValue := "plugins/" + pluginName + "/" + goPackage
						_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{tagValue})
					} else {
						for _, group := range parse.Service.Groups {
							for _, route := range group.Routes {
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
			} else if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
				if goPackage != "" {
					_ = g.Set(fmt.Sprintf("paths.%s.%s.tags", pmk, pmmk), []string{goPackage})
				} else {
					for _, group := range parse.Service.Groups {
						for _, route := range group.Routes {
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

			pluginName = getPluginNameFromFilePath(apiPath)
			if pluginName != "" {
				operationID := cast.ToString(g.Get(fmt.Sprintf("paths.%s.%s.operationId", pmk, pmmk)))
				if operationID != "" {
					_ = g.Set(fmt.Sprintf("paths.%s.%s.operationId", pmk, pmmk), "plugins/"+pluginName+"/"+operationID)
				}
			}

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

	return os.WriteFile(filepath.Join(outputDir, apiFile+".json"), encodeToJSON, 0o644)
}

func runProtoSwagger(files []string, progressChan chan<- progress.Message) error {
	if err := ensureSwaggerOutputDir(); err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	var eg errgroup.Group
	eg.SetLimit(len(files))
	for _, protoPath := range files {
		protoPath := protoPath
		eg.Go(func() error {
			if err := processSwaggerProtoFile(protoPath); err != nil {
				return errors.Wrapf(err, "swagger proto file: %s", protoPath)
			}
			if progressChan != nil {
				progressChan <- progress.NewFile(protoPath)
			}
			return nil
		})
	}

	return eg.Wait()
}

func listSwaggerProtoFiles() ([]string, error) {
	var files []string

	switch {
	case len(config.C.Gen.Swagger.Desc) > 0:
		for _, v := range config.C.Gen.Swagger.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					files = append(files, v)
				}
				continue
			}

			specifiedProtoFiles, err := desc.FindRpcServiceProtoFiles(v)
			if err != nil {
				return nil, err
			}
			files = append(files, specifiedProtoFiles...)
		}
	default:
		if pathx.FileExists(config.C.ProtoDir()) {
			var err error
			files, err = desc.FindRpcServiceProtoFiles(config.C.ProtoDir())
			if err != nil {
				return nil, err
			}
		}

		plugins, err := serverless.GetPlugins()
		if err == nil {
			for _, p := range plugins {
				if pathx.FileExists(filepath.Join(p.Path, "desc", "proto")) {
					pluginFiles, err := desc.FindRpcServiceProtoFiles(filepath.Join(p.Path, "desc", "proto"))
					if err != nil {
						return nil, err
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
			continue
		}

		specifiedProtoFiles, err := desc.FindRpcServiceProtoFiles(v)
		if err != nil {
			return nil, err
		}
		for _, saf := range specifiedProtoFiles {
			files = lo.Reject(files, func(item string, _ int) bool {
				return item == saf
			})
		}
	}

	return files, nil
}

func processSwaggerProtoFile(protoPath string) error {
	pluginName := getPluginNameFromFilePath(protoPath)
	var pluginProtoDir string
	var outputDir string

	if pluginName != "" {
		protoDirPrefix := filepath.Join("", config.C.ProtoDir()) + string(filepath.Separator)
		descProtoIndex := strings.Index(protoPath, protoDirPrefix)
		if descProtoIndex == -1 {
			if strings.HasSuffix(filepath.Dir(protoPath), filepath.Join("", config.C.ProtoDir())) {
				pluginProtoDir = filepath.Dir(protoPath)
			} else {
				return fmt.Errorf("invalid plugin proto path: %s", protoPath)
			}
		} else {
			pluginProtoDir = protoPath[:descProtoIndex+len(protoDirPrefix)]
		}
		outputDir = filepath.Join(config.C.Gen.Swagger.Output, "plugins", pluginName)
	} else {
		outputDir = config.C.Gen.Swagger.Output
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	var includeArgs []string
	if pluginName != "" {
		includeArgs = append(includeArgs, "-I"+pluginProtoDir)
		pluginThirdParty := filepath.Join(pluginProtoDir, "third_party")
		if pathx.FileExists(pluginThirdParty) {
			includeArgs = append(includeArgs, "-I"+pluginThirdParty)
		}
	}
	includeArgs = append(includeArgs, "-I"+config.C.ProtoDir())
	includeArgs = append(includeArgs, "-I"+filepath.Join(config.C.ProtoDir(), "third_party"))

	command := fmt.Sprintf("protoc %s %s --openapiv2_out=%s",
		strings.Join(includeArgs, " "),
		protoPath,
		outputDir,
	)
	_, err := execx.Run(command, config.C.Wd())
	return err
}

func runMergeSwagger(progressChan chan<- progress.Message) error {
	outputFile := filepath.Join(config.C.Gen.Swagger.Output, "swagger.json")
	if err := mergeSwaggerFiles(); err != nil {
		return errors.Wrapf(err, "merge swagger file: %s", outputFile)
	}
	if progressChan != nil {
		progressChan <- progress.NewFile(outputFile)
	}
	return nil
}

func ensureSwaggerOutputDir() error {
	return os.MkdirAll(config.C.Gen.Swagger.Output, 0o755)
}

func hasSwaggerSourceInput() bool {
	if pathx.FileExists(config.C.ApiDir()) {
		if files, err := desc.FindApiFiles(config.C.ApiDir()); err == nil && len(files) > 0 {
			return true
		}
	}

	if pathx.FileExists(config.C.ProtoDir()) {
		if files, err := desc.FindExcludeThirdPartyProtoFiles(config.C.ProtoDir()); err == nil && len(files) > 0 {
			return true
		}
	}

	plugins, err := serverless.GetPlugins()
	if err != nil {
		return false
	}

	for _, p := range plugins {
		pluginAPIDir := filepath.Join(p.Path, "desc", "api")
		if pathx.FileExists(pluginAPIDir) {
			if files, err := desc.FindApiFiles(pluginAPIDir); err == nil && len(files) > 0 {
				return true
			}
		}

		pluginProtoDir := filepath.Join(p.Path, "desc", "proto")
		if pathx.FileExists(pluginProtoDir) {
			if files, err := desc.FindExcludeThirdPartyProtoFiles(pluginProtoDir); err == nil && len(files) > 0 {
				return true
			}
		}
	}

	return false
}

func splitPluginPaths(files []string) (regular []string, plugin []string) {
	for _, file := range files {
		if getPluginNameFromFilePath(file) != "" {
			plugin = append(plugin, file)
			continue
		}
		regular = append(regular, file)
	}

	return regular, plugin
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
			// Error removed("failed to read swagger file %s: %v", filePath, err)
			continue
		}

		g, err := genius.NewFromRawJSON(file)
		if err != nil {
			// Error removed("failed to parse swagger file %s: %v", filePath, err)
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
