package desc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

const (
	metadataDir    = ".jzero"
	metadataSubDir = "desc-metadata"
	metadataFile   = "metadata.json"
)

type Metadata struct {
	API   *APIMetadata   `json:"api,omitempty"`
	Proto *ProtoMetadata `json:"proto,omitempty"`
}

type APIMetadata struct {
	Routes []APIRoute `json:"routes"`
}

type APIRoute struct {
	Handler  string `json:"handler"`
	Group    string `json:"group"`
	Path     string `json:"path"`
	Logic    string `json:"logic"`
	Desc     string `json:"desc"`
	DescLine int    `json:"desc-line"`
}

type ProtoMetadata struct {
	RPCs []ProtoRPC `json:"rpcs"`
}

type ProtoRPC struct {
	RPC      string `json:"rpc"`
	Service  string `json:"service"`
	Path     string `json:"path"`
	Logic    string `json:"logic"`
	Desc     string `json:"desc"`
	DescLine int    `json:"desc-line"`
}

// CollectFromAPI 从 API 文件收集元数据
// apiSpecMap: 已解析的 API 文件映射，key 为文件路径
func CollectFromAPI(apiSpecMap map[string]*spec.ApiSpec) (*APIMetadata, error) {
	var metadata APIMetadata

	// 为每个 API 文件解析行号信息
	lineMapCache := make(map[string]map[string]int)
	for apiFile := range apiSpecMap {
		lineMap, err := parseAPILineNumbers(apiFile)
		if err != nil {
			logx.Debugf("parse api line numbers failed: %s", err.Error())
		} else {
			lineMapCache[apiFile] = lineMap
		}
	}

	for apiFile, apiSpec := range apiSpecMap {
		for _, group := range apiSpec.Service.Groups {
			groupAnnotation := group.GetAnnotation("group")

			for _, route := range group.Routes {
				handler := route.Handler
				handler = strings.TrimSuffix(route.Handler, "Handler")

				// 构建 logic 文件路径
				logicPath, err := buildLogicPath(groupAnnotation, handler)
				if err != nil {
					logx.Debugf("build logic path failed: %s", err.Error())
					continue
				}

				// 获取行号
				var lineNumber int
				if lineMap, ok := lineMapCache[apiFile]; ok {
					if line, ok := lineMap[handler]; ok {
						lineNumber = line
					}
				}

				apiRoute := APIRoute{
					Handler:  handler,
					Group:    groupAnnotation,
					Path:     fmt.Sprintf("%s:%s", strings.ToLower(route.Method), filepath.ToSlash(filepath.Join(group.GetAnnotation("prefix"), route.Path))),
					Logic:    logicPath,
					Desc:     filepath.ToSlash(filepath.Join(config.C.Wd(), apiFile)),
					DescLine: lineNumber,
				}

				metadata.Routes = append(metadata.Routes, apiRoute)
			}
		}
	}

	return &metadata, nil
}

// CollectFromProto 从 Proto 文件收集元数据
// protoSpecMap: 已解析的 Proto 文件映射，key 为文件路径
func CollectFromProto(protoSpecMap map[string]rpcparser.Proto) (*ProtoMetadata, error) {
	var metadata ProtoMetadata

	// 为每个 Proto 文件解析行号信息
	lineMapCache := make(map[string]map[string]int)
	for protoFile := range protoSpecMap {
		lineMap, err := parseProtoLineNumbers(protoFile, protoSpecMap[protoFile])
		if err != nil {
			logx.Debugf("parse proto line numbers failed: %s", err.Error())
		} else {
			lineMapCache[protoFile] = lineMap
		}
	}

	for protoFile, protoSpec := range protoSpecMap {
		for _, service := range protoSpec.Service {
			for _, rpc := range service.RPC {
				logicPath, err := buildLogicPath(strings.ToLower(service.Name), rpc.Name)
				if err != nil {
					continue
				}

				// 获取行号
				var lineNumber int
				if lineMap, ok := lineMapCache[protoFile]; ok {
					key := service.Name + "." + rpc.Name
					if line, ok := lineMap[key]; ok {
						lineNumber = line
					}
				}

				protoRPC := ProtoRPC{
					RPC:      rpc.Name,
					Service:  service.Name,
					Path:     fmt.Sprintf("/%s.%s/%s", protoSpec.Package.Name, service.Name, rpc.Name),
					Logic:    logicPath,
					Desc:     filepath.ToSlash(filepath.Join(config.C.Wd(), protoFile)),
					DescLine: lineNumber,
				}

				metadata.RPCs = append(metadata.RPCs, protoRPC)
			}
		}
	}

	return &metadata, nil
}

// Save 保存元数据到文件
func Save(metadata *Metadata) error {
	projectPath, err := getProjectPath()
	if err != nil {
		return err
	}

	// 构建元数据目录路径
	metadataPath := filepath.Join(getHomeDir(), metadataDir, metadataSubDir, projectPath)

	// 创建目录
	if err = os.MkdirAll(metadataPath, 0o755); err != nil {
		return errors.Wrapf(err, "create metadata dir: %s", metadataPath)
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal metadata")
	}

	// 写入文件
	metadataFilepath := filepath.Join(metadataPath, metadataFile)
	if err := os.WriteFile(metadataFilepath, data, 0o644); err != nil {
		return errors.Wrapf(err, "write metadata file: %s", metadataFilepath)
	}

	return nil
}

// buildLogicPath 构建 logic 文件路径（绝对路径）
func buildLogicPath(groupOrService, handler string) (string, error) {
	wd := config.C.Wd()
	var logicPath string
	formated, err := format.FileNamingFormat(config.C.Style, handler)
	if err != nil {
		return "", err
	}
	if groupOrService != "" {
		logicPath = filepath.Join(wd, "internal", "logic", groupOrService, formated+".go")
	} else {
		logicPath = filepath.Join(wd, "internal", "logic", formated+".go")
	}

	return filepath.ToSlash(logicPath), nil
}

// parseProtoLineNumbers 解析 Proto 文件并提取 RPC 行号
func parseProtoLineNumbers(protoFile string, protoSpec rpcparser.Proto) (map[string]int, error) {
	lineMap := make(map[string]int)

	// 遍历服务和 RPC
	for _, service := range protoSpec.Service {
		for _, rpc := range service.RPC {
			// 使用 rpc.Position.Line 获取 RPC 定义所在的行号
			// rpcparser.RPC 嵌入了 *proto.RPC，它有 Position 字段
			if rpc.Position.Line > 0 {
				key := service.Name + "." + rpc.Name
				lineMap[key] = rpc.Position.Line
			}
		}
	}

	return lineMap, nil
}

// getProjectPath 获取项目路径（相对于 home 的唯一标识）
func getProjectPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 使用完整路径作为项目标识符，去掉 / 前缀
	projectPath := strings.TrimPrefix(filepath.ToSlash(wd), "/")
	return projectPath, nil
}

// getHomeDir 获取用户主目录
func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果获取失败，使用当前目录
		homeDir = "."
	}
	return homeDir
}

// parseAPILineNumbers 解析 API 文件并提取路由行号
func parseAPILineNumbers(apiFile string) (map[string]int, error) {
	p := parser.New(apiFile, "")
	astAST := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return nil, err
	}

	if astAST == nil {
		return nil, nil
	}

	lineMap := make(map[string]int)

	for _, stmt := range astAST.Stmts {
		serviceStmt, ok := stmt.(*ast.ServiceStmt)
		if !ok {
			continue
		}

		for _, r := range serviceStmt.Routes {
			lineMap[r.AtHandler.Name.RawText()] = r.AtHandler.AtHandler.Pos().Line
		}
	}
	return lineMap, nil
}
