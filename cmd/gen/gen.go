package gen

import (
	"fmt"
	"github.com/jaronnie/genius"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jaronnie/jzero/daemon/pkg/mod"
	"github.com/jaronnie/jzero/daemon/pkg/stringx"
	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"
)

var (
	WorkingDir string
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

	configBytes, err := os.ReadFile(filepath.Join(wd, "config.toml"))
	cobra.CheckErr(err)

	g, err := genius.NewFromToml(configBytes)
	cobra.CheckErr(err)

	// read proto dir
	protoDir, err := os.ReadDir(filepath.Join(wd, "daemon", "desc", "proto"))
	if len(protoDir) == 0 {
		fmt.Printf("proto dir [%s] not found. Skip generate", filepath.Join(wd, "daemon", "desc", "proto"))
	}

	// api file path
	apiFilePath := filepath.Join(wd, "daemon", "desc", "api", cast.ToString(g.Get("APP"))+".api")

	err = check(g, wd, protoDir, apiFilePath)
	cobra.CheckErr(err)

	var protosets []string
	var serverImports ImportLines
	var pbImports ImportLines
	var registerServers RegisterLines

	moduleStruct, err := mod.GetGoMod(wd)
	cobra.CheckErr(err)

	// 正常删除无用文件夹
	defer func() {
		_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", cast.ToString(g.Get("APP")))))
		_ = os.RemoveAll(filepath.Join(wd, "daemon", "etc"))
		os.Exit(-1)
	}()

	// 异常删除无用文件夹
	go extraFileHandler(g, wd)

	for _, v := range protoDir {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "proto") {
			fmt.Printf("%s to generate proto code. \n%s proto file %s\n", color.WithColor("Start", color.FgGreen), color.WithColor("Using", color.FgGreen), filepath.Join(wd, "daemon", "desc", "proto", v.Name()))
			command := fmt.Sprintf("goctl rpc protoc daemon/desc/proto/%s  -I./daemon/desc/proto --go_out=./daemon/internal --go-grpc_out=./daemon/internal  --zrpc_out=./daemon --client=false --home %s -m", v.Name(), filepath.Join(wd, ".template", "go-zero"))
			_, err := execx.Run(command, wd)
			cobra.CheckErr(err)
			fmt.Println(color.WithColor("Done", color.FgGreen))

			fileBase := v.Name()[0 : len(v.Name())-len(path.Ext(v.Name()))]
			_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", fileBase)))

			// # gen proto descriptor
			_ = os.MkdirAll(filepath.Join(wd, ".protosets"), 0o755)
			protocCommand := fmt.Sprintf("protoc --include_imports -I./daemon/desc/proto --descriptor_set_out=.protosets/%s.pb daemon/desc/proto/%s.proto", fileBase, fileBase)
			_, err = execx.Run(protocCommand, wd)
			cobra.CheckErr(err)

			protosets = append(protosets, filepath.Join(".protosets", fmt.Sprintf("%s.pb", fileBase)))

			// parse proto
			protoParser := rpcparser.NewDefaultProtoParser()
			parse, err := protoParser.Parse(filepath.Join(wd, "daemon", "desc", "proto", v.Name()), true)
			cobra.CheckErr(err)
			for _, s := range parse.Service {
				serverImports = append(serverImports, fmt.Sprintf(`%ssvr "%s/daemon/internal/server/%s"`, s.Name, moduleStruct.Path, s.Name))
				registerServers = append(registerServers, fmt.Sprintf("%s.Register%sServer(grpcServer, %ssvr.New%sServer(ctx))", filepath.Base(parse.GoPackage), stringx.FirstUpper(s.Name), s.Name, stringx.FirstUpper(s.Name)))
			}
			pbImports = append(pbImports, fmt.Sprintf(`"%s/daemon/internal/%s"`, moduleStruct.Path, strings.TrimPrefix(parse.GoPackage, "./")))
		}
	}

	// 生成 api 代码
	if pathx.FileExists(apiFilePath) {
		fmt.Printf("%s to generate api code.\n%s api file %s\n", color.WithColor("Start", color.FgGreen), color.WithColor("Using", color.FgGreen), apiFilePath)
		command := fmt.Sprintf("goctl api go --api %s --dir ./daemon --home %s", apiFilePath, filepath.Join(wd, ".template", "go-zero"))
		_, err = execx.Run(command, wd)
		cobra.CheckErr(err)
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	// 检测是否包含 sql
	sqlDir := filepath.Join(wd, "daemon", "desc", "sql")
	if f, err := os.Stat(sqlDir); err == nil && f.IsDir() {
		fs, err := os.ReadDir(sqlDir)
		cobra.CheckErr(err)
		for _, f := range fs {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
				sqlFilePath := filepath.Join(sqlDir, f.Name())
				fmt.Printf("%s to generate model code.\n%s sql file %s\n", color.WithColor("Start", color.FgGreen), color.WithColor("Using", color.FgGreen), sqlFilePath)
				command := fmt.Sprintf("goctl model mysql ddl --src daemon/desc/sql/%s --dir ./daemon/internal/model/%s --home %s", f.Name(), f.Name()[0:len(f.Name())-len(path.Ext(f.Name()))], filepath.Join(wd, ".template", "go-zero"))
				_, err = execx.Run(command, wd)
				cobra.CheckErr(err)
				fmt.Println(color.WithColor("Done", color.FgGreen))
			}
		}
	}

	// 生成 daemon/zrpc.go
	fmt.Printf("%s to generate daemon/zrpc.go\n", color.WithColor("Start", color.FgGreen))
	zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module":          moduleStruct.Path,
		"APP":             cast.ToString(g.Get("APP")),
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "zrpc.go.tpl")))
	err = os.WriteFile(filepath.Join(wd, "daemon", "zrpc.go"), zrpcFile, 0o644)
	cobra.CheckErr(err)
	fmt.Printf("%s", color.WithColor("Done\n", color.FgGreen))

	// 修改 config.toml protosets 内容
	// 检测是否需要修改 config.toml. 以及让用户选择是否自动更新文件
	existProtosets := g.Get("Gateway.Upstreams.0.ProtoSets")
	if len(lo.Intersect(cast.ToStringSlice(existProtosets), protosets)) != len(protosets) {
		var in string
		fmt.Println("检测到 config.toml 中 Gateway.Upstreams.0.ProtoSets 配置需要更新. 是否自动更新 y/n. 更新需谨慎, 会将注释删掉")
		_, _ = fmt.Scanln(&in)
		switch {
		case strings.EqualFold(in, "y"):
			fmt.Printf("%s to update config.toml\n", color.WithColor("Start", color.FgGreen))
			err = g.Set("Gateway.Upstreams.0.ProtoSets", protosets)
			cobra.CheckErr(err)
			toml, err := g.EncodeToToml()
			cobra.CheckErr(err)
			err = os.WriteFile(filepath.Join(wd, "config.toml"), toml, 0644)
			cobra.CheckErr(err)
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		case strings.EqualFold(in, "n"):
			fmt.Printf("请手动更新 Gateway.Upstreams.0.ProtoSets 配置\n配置该值为: \n%s\n",
				color.WithColor(fmt.Sprintf("%v", protosets), color.FgGreen))
		}
	}

	return nil
}

func check(g *genius.Genius, wd string, protoDir []os.DirEntry, apiFilePath string) error {
	// check logic dir duplicate
	var protoLogicDir []string
	var apiLogicDir []string

	for _, v := range protoDir {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "proto") {
			protoParser := rpcparser.NewDefaultProtoParser()
			parse, err := protoParser.Parse(filepath.Join(wd, "daemon", "desc", "proto", v.Name()), true)
			cobra.CheckErr(err)
			for _, s := range parse.Service {
				protoLogicDir = append(protoLogicDir, s.Name)
			}
		}
	}

	// parse api
	if pathx.FileExists(apiFilePath) {
		apiSpec, err := parser.Parse(filepath.Join(wd, "daemon", "desc", "api", cast.ToString(g.Get("APP"))+".api"))
		if err != nil {
			return err
		}
		for _, v := range apiSpec.Service.Groups {
			if g, ok := v.Annotation.Properties["group"]; ok {
				apiLogicDir = append(apiLogicDir, g)
			}
		}
	}

	// check logicDir has duplicate
	intersect := lo.Intersect(protoLogicDir, apiLogicDir)
	if len(intersect) != 0 {
		errorString := fmt.Sprintf("service %s duplicate. Please check api file group and proto file service name", strings.Join(intersect, ","))
		return errors.New(errorString)
	}

	return nil
}

func extraFileHandler(g *genius.Genius, wd string) {
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", cast.ToString(g.Get("APP")))))
			_ = os.RemoveAll(filepath.Join(wd, "daemon", "etc"))
			os.Exit(-1)
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func init() {
	logx.Disable()
}
