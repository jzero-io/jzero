package desc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

func GetFrameType() (string, error) {
	// 判断 core 项目类型 api/rpc
	var frameType string
	if _, err := os.Stat(filepath.Join("desc", "api")); err == nil {
		// api 项目
		frameType = "api"
	}
	if _, err := os.Stat(filepath.Join("desc", "proto")); err == nil {
		// rpc 项目
		frameType = "rpc"

		// 获取全量 proto 文件
		protoFiles, err := GetProtoFilepath(config.C.ProtoDir())
		if err != nil {
			return "", err
		}

		for _, v := range protoFiles {
			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			var parse rpcparser.Proto
			parse, err = protoParser.Parse(v, true)
			if err != nil {
				return "", err
			}
			if IsNeedGenProtoDescriptor(parse) {
				frameType = "gateway"
				break
			}
		}
	}

	return frameType, nil
}

func GetProtoDescriptorPath(protoPath string) string {
	rel, err := filepath.Rel(filepath.Join("desc", "proto"), protoPath)
	if err != nil {
		return ""
	}

	return filepath.Join("desc", "pb", strings.TrimSuffix(rel, ".proto")+".pb")
}

func IsNeedGenProtoDescriptor(proto rpcparser.Proto) bool {
	for _, ps := range proto.Service {
		for _, rpc := range ps.RPC {
			for _, option := range rpc.Options {
				if option.Name == "(google.api.http)" {
					return true
				}
			}
		}
	}
	return false
}

func GetProtoFilepath(protoDirPath string) ([]string, error) {
	var protoFilenames []string

	protoDir, err := os.ReadDir(protoDirPath)
	if err != nil {
		return nil, nil
	}

	for _, protoFile := range protoDir {
		if protoFile.IsDir() {
			filenames, err := GetProtoFilepath(filepath.Join(protoDirPath, protoFile.Name()))
			if err != nil {
				return nil, err
			}
			protoFilenames = append(protoFilenames, filenames...)
		} else {
			if strings.HasSuffix(protoFile.Name(), ".proto") {
				if b, err := protoHasService(filepath.Join(protoDirPath, protoFile.Name())); err == nil && b {
					protoFilenames = append(protoFilenames, filepath.Join(protoDirPath, protoFile.Name()))
				} else if err != nil {
					return nil, err
				}
			}
		}
	}
	return protoFilenames, nil
}

func protoHasService(fp string) (bool, error) {
	r := rpcparser.DefaultProtoParser{}

	parse, err := r.Parse(fp, true)
	if err != nil {
		if strings.Contains(err.Error(), "rpc service not found") {
			return false, nil
		}
		return false, err
	}
	return len(parse.Service) > 0, nil
}

func GetApiServiceName(apiDirName string, files ...string) string {
	if len(files) == 0 {
		fs, err := getApiFileRelPath(apiDirName)
		if err != nil {
			return ""
		}
		for _, file := range fs {
			apiSpec, err := parser.Parse(filepath.Join(apiDirName, file), "")
			if err != nil {
				return ""
			}
			if apiSpec.Service.Name != "" {
				return apiSpec.Service.Name
			}
		}
	} else {
		file := files[0]
		apiSpec, err := parser.Parse(file, "")
		if err != nil {
			return ""
		}
		if apiSpec.Service.Name != "" {
			return apiSpec.Service.Name
		}
	}

	return ""
}

func GetRpcMethodUrl(method *descriptorpb.MethodDescriptorProto) string {
	ext := proto.GetExtension(method.GetOptions(), annotations.E_Http)
	switch rule := ext.(type) {
	case *annotations.HttpRule:
		switch httpRule := rule.GetPattern().(type) {
		case *annotations.HttpRule_Get:
			return "GET:" + httpRule.Get
		case *annotations.HttpRule_Post:
			return "POST:" + httpRule.Post
		case *annotations.HttpRule_Put:
			return "PUT:" + httpRule.Put
		case *annotations.HttpRule_Delete:
			return "DELETE:" + httpRule.Delete
		case *annotations.HttpRule_Patch:
			return "PATCH:" + httpRule.Patch
		}
	}
	return ""
}

func getApiFileRelPath(apiDirName string) ([]string, error) {
	var apiFiles []string

	allApiFiles, err := FindApiFiles(apiDirName)
	if err != nil {
		return nil, err
	}
	for _, file := range allApiFiles {
		rel, err := filepath.Rel(apiDirName, file)
		if err != nil {
			return nil, err
		}
		apiFiles = append(apiFiles, filepath.ToSlash(rel))
	}

	return apiFiles, nil
}

func findDescFiles(dir, descExt string) ([]string, error) {
	var descFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			subFiles, err := findDescFiles(filepath.Join(dir, file.Name()), descExt)
			if err != nil {
				return nil, err
			}
			descFiles = append(descFiles, subFiles...)
		} else if filepath.Ext(file.Name()) == descExt {
			descFiles = append(descFiles, filepath.Join(dir, file.Name()))
		}
	}

	return descFiles, nil
}

func FindApiFiles(dir string) ([]string, error) {
	return findDescFiles(dir, ".api")
}

func FindSqlFiles(dir string) ([]string, error) {
	return findDescFiles(dir, ".sql")
}

func FindRouteApiFiles(dir string) ([]string, error) {
	var routeFiles []string
	files, err := findDescFiles(dir, ".api")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		parse, err := parser.Parse(f, "")
		if err != nil {
			return nil, errors.Errorf("parse api file: %s", f)
		}
		if len(parse.Service.Routes()) > 0 {
			routeFiles = append(routeFiles, f)
		}
	}
	return routeFiles, nil
}

// GetApiFrameMainGoFilename goctl/api/gogen/genmain.go
func GetApiFrameMainGoFilename(wd, file, style string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"), file)
	serviceName = strings.ToLower(serviceName)
	filename, err := format.FileNamingFormat(style, serviceName)
	if err != nil {
		return ""
	}

	if strings.HasSuffix(filename, "-api") {
		filename = strings.ReplaceAll(filename, "-api", "")
	}
	return filename + ".go"
}

// GetApiFrameEtcFilename goctl/api/gogen/genetc.go
func GetApiFrameEtcFilename(wd, file, style string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"), file)
	filename, err := format.FileNamingFormat(style, serviceName)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}

// GetProtoFrameMainGoFilename goctl/rpc/generator/genmain.go
func GetProtoFrameMainGoFilename(source, style string) string {
	filename, err := format.FileNamingFormat(style, source)
	if err != nil {
		return ""
	}
	return filename + ".go"
}

// GetProtoFrameEtcFilename goctl/rpc/generator/genetc.go
func GetProtoFrameEtcFilename(source, style string) string {
	filename, err := format.FileNamingFormat(style, source)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}
