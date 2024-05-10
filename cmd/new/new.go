package new

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaronnie/genius"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"

	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/jzero-io/jzero/embeded"
)

var (
	Module string
	Dir    string
	APP    string
	// ConfigType config type
	ConfigType string

	Version string
)

func NewProject(_ *cobra.Command, _ []string) error {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if embeded.Home == "" {
		embeded.Home = filepath.Join(homeDir, ".jzero", Version)
	}

	// mkdir output
	err = os.MkdirAll(Dir, 0o755)
	cobra.CheckErr(err)
	// go mod init
	_, err = execx.Run(fmt.Sprintf("go mod init %s", Module), Dir)
	cobra.CheckErr(err)
	// touch main.go
	mainFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "main.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "main.go"), mainFile)
	cobra.CheckErr(err)
	// mkdir cmd dir
	err = os.MkdirAll(filepath.Join(Dir, "cmd"), 0o755)
	cobra.CheckErr(err)
	// touch cmd/root.go
	rootCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module":     Module,
		"APP":        APP,
		"ConfigType": ConfigType,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "root.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "cmd", "root.go"), rootCmdFile)
	cobra.CheckErr(err)
	// touch cmd/server.go
	serverCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "server.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "cmd", "server.go"), serverCmdFile)
	cobra.CheckErr(err)
	// touch cmd/version.go
	versionCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "version.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "cmd", "version.go"), versionCmdFile)
	cobra.CheckErr(err)
	// mkdir app dir
	err = os.MkdirAll(filepath.Join(Dir, "app"), 0o755)
	cobra.CheckErr(err)
	// touch app/server.go
	serverFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "server.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "server.go"), serverFile)
	cobra.CheckErr(err)

	// touch app/zrpc.go
	zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "zrpc.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "zrpc.go"), zrpcFile)
	cobra.CheckErr(err)

	// mkdir api, proto dir
	// write proto dir
	if zrpcFile != nil {
		err = embeded.WriteTemplateDir(filepath.Join("jzero", "app", "desc", "proto"), filepath.Join(Dir, "app", "desc", "proto"))
		cobra.CheckErr(err)
	}

	err = os.MkdirAll(filepath.Join(Dir, "app", "desc", "api"), 0o755)
	cobra.CheckErr(err)
	// touch app/desc/api/{{.APP}}.api
	err = checkWrite(filepath.Join(Dir, "app", "desc", "api", APP+".api"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", "jzero.api.tpl")))
	cobra.CheckErr(err)
	// touch app/desc/api/hello.api
	helloApiFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", "hello.api.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "desc", "api", "hello.api"), helloApiFile)
	cobra.CheckErr(err)

	// write config.yaml
	configYamlFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "config.yaml.tpl")))
	cobra.CheckErr(err)

	g, err := genius.NewFromYaml(configYamlFile)
	cobra.CheckErr(err)

	switch ConfigType {
	case "toml":
		configTomlFile, err := g.EncodeToToml()
		cobra.CheckErr(err)
		err = checkWrite(filepath.Join(Dir, "config.toml"), configTomlFile)
		cobra.CheckErr(err)
	case "yaml":
		err = checkWrite(filepath.Join(Dir, "config.yaml"), configYamlFile)
		cobra.CheckErr(err)
	case "json":
		configJsonFile, err := g.EncodeToJSON()
		cobra.CheckErr(err)
		err = checkWrite(filepath.Join(Dir, "config.json"), configJsonFile)
		cobra.CheckErr(err)
	}

	// ################# start gen config ###################
	// write app/internal/config/config.go
	err = os.MkdirAll(filepath.Join(Dir, "app", "internal", "config"), 0o755)
	cobra.CheckErr(err)

	configGoFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "config", "config.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "internal", "config", "config.go"), configGoFile)
	cobra.CheckErr(err)
	// ################# end gen config ###################

	// ################# start gen middlewares ###################
	// write app/middlewares/response.go
	err = os.MkdirAll(filepath.Join(Dir, "app", "middlewares"), 0o755)
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "middlewares", "response.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "response.go.tpl")))
	cobra.CheckErr(err)

	// write app/middlewares/errors.go
	err = checkWrite(filepath.Join(Dir, "app", "middlewares", "errors.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "errors.go.tpl")))
	cobra.CheckErr(err)

	// write app/middlewares/grpc_rate_limit.go
	err = checkWrite(filepath.Join(Dir, "app", "middlewares", "grpc_rate_limit.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "grpc_rate_limit.go.tpl")))
	cobra.CheckErr(err)

	// write app/middlewares/logs.go
	logsMiddlewareFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":    APP,
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "logs.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "middlewares", "logs.go"), logsMiddlewareFile)
	cobra.CheckErr(err)

	// ################# end gen middlewares ###################

	// write app/internal/handler/myroutes.go
	_ = os.MkdirAll(filepath.Join(Dir, "app", "internal", "handler"), 0o755)
	myroutesFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "handler", "myroutes.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "internal", "handler", "myroutes.go"), myroutesFile)
	cobra.CheckErr(err)

	// write app/internal/handler/myhandler.go
	myhandlerFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "handler", "myhandler.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "internal", "handler", "myhandler.go"), myhandlerFile)
	cobra.CheckErr(err)

	// write Dockerfile
	dockerFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":        APP,
		"ConfigType": ConfigType,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "Dockerfile"), dockerFile)
	cobra.CheckErr(err)

	// write Dockerfile-arm64
	dockerArm64File, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":        APP,
		"ConfigType": ConfigType,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile-arm64.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "Dockerfile-arm64"), dockerArm64File)
	cobra.CheckErr(err)

	// write .goreleaser.yaml
	goreleaserBytes := embeded.ReadTemplateFile(filepath.Join("jzero", "goreleaser.yaml.tpl"))
	goreleaserBytes = bytes.ReplaceAll(goreleaserBytes, []byte("{{ .APP }}"), []byte(APP))
	goreleaserBytes = bytes.ReplaceAll(goreleaserBytes, []byte("{{ .Module }}"), []byte(Module))
	err = checkWrite(filepath.Join(Dir, ".goreleaser.yaml"), goreleaserBytes)
	cobra.CheckErr(err)

	// write Taskfile.yml
	err = checkWrite(filepath.Join(Dir, "Taskfile.yml"), embeded.ReadTemplateFile(filepath.Join("jzero", "Taskfile.yml.tpl")))
	cobra.CheckErr(err)

	err = embeded.WriteTemplateDir(filepath.Join("go-zero"), filepath.Join(embeded.Home, "go-zero"))
	cobra.CheckErr(err)

	// write .gitignore
	gitignoreFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "gitignore.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, ".gitignore"), gitignoreFile)
	cobra.CheckErr(err)

	return nil
}

func CheckWrite(path string, bytes []byte) error {
	return checkWrite(path, bytes)
}

func checkWrite(path string, bytes []byte) error {
	if bytes == nil || len(bytes) == 0 {
		return nil
	}
	err := os.WriteFile(path, bytes, 0o644)
	return err
}
