package gen

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jaronnie/genius"
	"github.com/jzero-io/jzero/app/pkg/mod"
	"github.com/jzero-io/jzero/app/pkg/stringx"
	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/jzero-io/jzero/embeded"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	WorkingDir string

	Version string
)

type (
	ImportLines   []string
	RegisterLines []string
)

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

func Gen(_ *cobra.Command, _ []string) error {
	// change dir
	if WorkingDir != "" {
		err := os.Chdir(WorkingDir)
		cobra.CheckErr(err)
	}

	wd, err := os.Getwd()
	cobra.CheckErr(err)
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), wd)

	if embeded.Home == "" {
		home, _ := os.UserHomeDir()
		embeded.Home = filepath.Join(home, ".jzero", Version)
	}

	// get configType
	configType, err := stringx.GetConfigType(wd)
	cobra.CheckErr(err)

	configBytes, err := os.ReadFile(filepath.Join(wd, "config."+configType))
	cobra.CheckErr(err)

	g, err := genius.NewFromType(configBytes, configType)
	cobra.CheckErr(err)

	// read proto dir
	protoDir, err := GetProtoDir(wd)
	cobra.CheckErr(err)

	var protosets []string
	var serverImports ImportLines
	var pbImports ImportLines
	var registerServers RegisterLines

	moduleStruct, err := mod.GetGoMod(wd)
	cobra.CheckErr(err)

	// 正常删除无用文件夹
	defer func() {
		removeExtraFiles(wd)
		os.Exit(0)
	}()

	// 异常删除无用文件夹
	go extraFileHandler(wd)

	for _, v := range protoDir {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "proto") {
			fmt.Printf("%s to generate proto code. \n%s proto file %s\n", color.WithColor("Start", color.FgGreen), color.WithColor("Using", color.FgGreen), filepath.Join(wd, "app", "desc", "proto", v.Name()))
			command := fmt.Sprintf("goctl rpc protoc app/desc/proto/%s  -I./app/desc/proto --go_out=./app/internal --go-grpc_out=./app/internal --zrpc_out=./app --client=false --home %s -m", v.Name(), filepath.Join(embeded.Home, "go-zero"))
			_, err := execx.Run(command, wd)
			cobra.CheckErr(err)
			fmt.Println(color.WithColor("Done", color.FgGreen))

			fileBase := v.Name()[0 : len(v.Name())-len(path.Ext(v.Name()))]
			_ = os.Remove(filepath.Join(wd, "app", fmt.Sprintf("%s.go", fileBase)))

			// # gen proto descriptor
			_ = os.MkdirAll(filepath.Join(wd, ".protosets"), 0o755)
			protocCommand := fmt.Sprintf("protoc --include_imports -I./app/desc/proto --descriptor_set_out=.protosets/%s.pb app/desc/proto/%s.proto", fileBase, fileBase)
			_, err = execx.Run(protocCommand, wd)
			cobra.CheckErr(err)

			protosets = append(protosets, filepath.Join(".protosets", fmt.Sprintf("%s.pb", fileBase)))

			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			parse, err := protoParser.Parse(filepath.Join(wd, "app", "desc", "proto", v.Name()), true)
			cobra.CheckErr(err)
			for _, s := range parse.Service {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/app/internal/server/%s"`, s.Name, moduleStruct.Path, s.Name))
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), s.Name, stringx.FirstUpper(s.Name)))
			}
			pbImports = append(pbImports, fmt.Sprintf(`"%s/app/internal/%s"`, moduleStruct.Path, strings.TrimPrefix(parse.GoPackage, "./")))
		}
	}

	// 生成 api 代码
	apiDirName := filepath.Join(wd, "app", "desc", "api")
	if pathx.FileExists(apiDirName) {
		fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))
		err = generateApiCode(wd, GetMainApiFilePath(apiDirName))
		cobra.CheckErr(err)
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	// 检测是否包含 sql
	sqlDir := filepath.Join(wd, "app", "desc", "sql")
	if f, err := os.Stat(sqlDir); err == nil && f.IsDir() {
		fs, err := os.ReadDir(sqlDir)
		cobra.CheckErr(err)
		fmt.Printf("%s to generate model code.\n", color.WithColor("Start", color.FgGreen))
		for _, f := range fs {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
				sqlFilePath := filepath.Join(sqlDir, f.Name())
				fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), sqlFilePath)
				command := fmt.Sprintf("goctl model mysql ddl --src app/desc/sql/%s --dir ./app/internal/model/%s --home %s", f.Name(), f.Name()[0:len(f.Name())-len(path.Ext(f.Name()))], filepath.Join(wd, ".template", "go-zero"))
				_, err = execx.Run(command, wd)
				cobra.CheckErr(err)
			}
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	// 生成 app/zrpc.go
	if pathx.FileExists(filepath.Join(wd, "app", "zrpc.go")) {
		fmt.Printf("%s to generate app/zrpc.go\n", color.WithColor("Start", color.FgGreen))
		zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
			"Module":          moduleStruct.Path,
			"APP":             cast.ToString(g.Get("APP")),
			"ServerImports":   serverImports,
			"PbImports":       pbImports,
			"RegisterServers": registerServers,
		}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "zrpc.go.tpl")))
		cobra.CheckErr(err)
		err = os.WriteFile(filepath.Join(wd, "app", "zrpc.go"), zrpcFile, 0o644)
		cobra.CheckErr(err)
		fmt.Printf("%s", color.WithColor("Done\n", color.FgGreen))

		if g.Get("Gateway") != nil {
			// 修改 config.toml protosets 内容
			// 检测是否需要修改 config.toml. 以及让用户选择是否自动更新文件
			existProtosets := g.Get("Gateway.Upstreams.0.ProtoSets")
			if len(lo.Intersect(cast.ToStringSlice(existProtosets), protosets)) != len(protosets) {
				var in string
				fmt.Printf("检测到 config.%s 中 Gateway.Upstreams.0.ProtoSets 配置需要更新. 是否自动更新 y/n. 更新需谨慎, 会将注释删掉\n", configType)
				_, _ = fmt.Scanln(&in)
				switch {
				case strings.EqualFold(in, "y"):
					fmt.Printf("%s to update config.%s\n", color.WithColor("Start", color.FgGreen), configType)
					err = g.Set("Gateway.Upstreams.0.ProtoSets", protosets)
					cobra.CheckErr(err)
					configBytes, err := g.EncodeToType(configType)
					cobra.CheckErr(err)
					err = os.WriteFile(filepath.Join(wd, "config."+configType), configBytes, 0o644)
					cobra.CheckErr(err)
					fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
				case strings.EqualFold(in, "n"):
					fmt.Printf("请手动更新 Gateway.Upstreams.0.ProtoSets 配置\n配置该值为: \n%s\n",
						color.WithColor(fmt.Sprintf("%v", protosets), color.FgGreen))
				}
			}
		}
	}
	return nil
}

func extraFileHandler(wd string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			removeExtraFiles(wd)
			os.Exit(-1)
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func getApiFilaPath(apiDirName string) []string {
	var apiFiles []string
	_ = filepath.Walk(apiDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".api" {
			spec, err := parser.Parse(path, nil)
			if err != nil {
				return err
			}
			if len(spec.Service.Routes()) > 0 {
				rel, err := filepath.Rel(apiDirName, path)
				if err != nil {
					return err
				}
				apiFiles = append(apiFiles, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	return apiFiles
}

func GetApiServiceName(apiDirName string) string {
	fs := getApiFilaPath(apiDirName)
	for _, file := range fs {
		apiSpec, err := parser.Parse(filepath.Join(apiDirName, file), "")
		if err != nil {
			cobra.CheckErr(err)
		}
		if apiSpec.Service.Name != "" {
			return apiSpec.Service.Name
		}
	}
	return ""
}

func generateApiCode(wd string, mainApiFilePath string) error {
	if mainApiFilePath == "" {
		return errors.New("empty mainApiFilePath")
	}
	defer os.Remove(mainApiFilePath)

	fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), mainApiFilePath)
	command := fmt.Sprintf("goctl api go --api %s --dir ./app --home %s", mainApiFilePath, filepath.Join(embeded.Home, "go-zero"))
	if _, err := execx.Run(command, wd); err != nil {
		return err
	}
	return nil
}

func GetMainApiFilePath(apiDirName string) string {
	apiDir, err := os.ReadDir(apiDirName)
	if err != nil {
		return ""
	}

	var mainApiFilePath string

	for _, file := range apiDir {
		if file.Name() == "main.api" {
			mainApiFilePath = filepath.Join(apiDirName, file.Name())
			break
		}
	}

	if mainApiFilePath == "" {
		apiFilePath := getApiFilaPath(apiDirName)
		sb := strings.Builder{}
		sb.WriteString("syntax = \"v1\"")
		sb.WriteString("\n")

		for _, api := range apiFilePath {
			sb.WriteString(fmt.Sprintf("import \"%s\"\n", api))
		}

		f, err := os.CreateTemp(apiDirName, "*.api")
		if err != nil {
			return ""
		}

		_, err = f.WriteString(sb.String())
		if err != nil {
			return ""
		}
		mainApiFilePath = f.Name()
		f.Close()
	}
	return mainApiFilePath
}

func removeExtraFiles(wd string) {
	_ = os.RemoveAll(filepath.Join(wd, "app", "etc"))
	_ = os.Remove(filepath.Join(wd, "app", fmt.Sprintf("%s.go", GetApiServiceName(filepath.Join(wd, "app", "desc", "api")))))
}

func init() {
	logx.Disable()
}
