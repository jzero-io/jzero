package desc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func GetFrameType() string {
	// Check for rpc frame
	if _, err := os.Stat("desc/proto"); err == nil {
		return "rpc"
	}

	// Check for api frame
	if _, err := os.Stat("desc/api"); err == nil {
		apiFiles, err := FindApiFiles("desc/api")
		if err == nil && len(apiFiles) > 0 {
			return "api"
		}
	}

	return ""
}

func getProtoDir(protoDirPath string) ([]os.DirEntry, error) {
	protoDir, err := os.ReadDir(protoDirPath)
	if err != nil {
		return nil, nil
	}
	return protoDir, nil
}

func GetProtoFilepath(protoDirPath string) ([]string, error) {
	var protoFilenames []string

	protoDir, err := getProtoDir(protoDirPath)
	if err != nil {
		return nil, err
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

// CheckToolVersion checks if the required tool version matches the minimum requirement
func CheckToolVersion(toolName, minVersionStr string) error {
	// Get tool version
	var cmd string
	switch toolName {
	case "goctl":
		cmd = "goctl -v"
	case "protoc":
		cmd = "protoc --version"
	case "protoc-gen-openapiv2":
		cmd = "protoc-gen-openapiv2 --version"
	case "protoc-gen-doc":
		cmd = "protoc-gen-doc --version"
	default:
		return errors.Errorf("unsupported tool: %s", toolName)
	}

	resp, err := execx.Run(cmd, "")
	if err != nil {
		return errors.Errorf("failed to get %s version: %v", toolName, err)
	}

	// Parse version from response
	var currentVersionStr string
	parts := strings.Split(resp, " ")

	switch toolName {
	case "goctl":
		if len(parts) >= 3 {
			currentVersionStr = parts[2]
		}
	case "protoc":
		if len(parts) >= 2 {
			currentVersionStr = strings.TrimSpace(parts[1])
		}
	default:
		if len(parts) >= 1 {
			currentVersionStr = strings.TrimSpace(parts[0])
		}
	}

	if currentVersionStr == "" {
		return errors.Errorf("failed to parse %s version from: %s", toolName, resp)
	}

	// Compare versions
	currentVersion, err := version.NewVersion(currentVersionStr)
	if err != nil {
		return errors.Errorf("invalid current version %s for %s: %v", currentVersionStr, toolName, err)
	}

	minVersion, err := version.NewVersion(minVersionStr)
	if err != nil {
		return errors.Errorf("invalid minimum version %s for %s: %v", minVersionStr, toolName, err)
	}

	if currentVersion.LessThan(minVersion) {
		return errors.Errorf("%s version %s is less than required version %s", toolName, currentVersionStr, minVersionStr)
	}

	return nil
}

// CheckFrameToolVersions checks version compatibility for frame-specific tools
func CheckFrameToolVersions(frameType string) error {
	// Common tool version requirements
	if err := CheckToolVersion("goctl", "1.7.0"); err != nil {
		return err
	}

	// Frame-specific tool requirements
	switch frameType {
	case "rpc":
		if err := CheckToolVersion("protoc", "3.19.0"); err != nil {
			return err
		}
		// Note: protoc-gen-openapiv2 and protoc-gen-doc might not have --version flags
		// so we skip version check for now, just check if they exist
	case "api":
		// API frame only needs goctl
	}

	return nil
}
