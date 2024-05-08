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
	err = os.WriteFile(filepath.Join(Dir, "main.go"), mainFile, 0o644)
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
	err = os.WriteFile(filepath.Join(Dir, "cmd", "root.go"), rootCmdFile, 0o644)
	cobra.CheckErr(err)
	// touch cmd/app.go
	appCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "server.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "cmd", "server.go"), appCmdFile, 0o644)
	cobra.CheckErr(err)
	// touch cmd/version.go
	versionCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "version.go.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "cmd", "version.go"), versionCmdFile, 0o644)
	cobra.CheckErr(err)
	// mkdir app dir
	err = os.MkdirAll(filepath.Join(Dir, "app"), 0o755)
	cobra.CheckErr(err)
	// touch app/app.go
	appFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "app.go.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "server.go"), appFile, 0o644)
	cobra.CheckErr(err)

	// touch app/zrpc.go
	zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "zrpc.go.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "zrpc.go"), zrpcFile, 0o644)
	cobra.CheckErr(err)

	// mkdir api, proto dir
	err = os.MkdirAll(filepath.Join(Dir, "app", "desc", "proto"), 0o755)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Join(Dir, "app", "desc", "api"), 0o755)
	cobra.CheckErr(err)
	// touch app/desc/api/{{.APP}}.api
	err = os.WriteFile(filepath.Join(Dir, "app", "desc", "api", APP+".api"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", "jzero.api.tpl")), 0o644)
	cobra.CheckErr(err)
	// touch app/desc/api/hello.api
	helloApiFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", "hello.api.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "desc", "api", "hello.api"), helloApiFile, 0o644)
	cobra.CheckErr(err)

	// write proto dir
	err = embeded.WriteTemplateDir(filepath.Join("jzero", "app", "desc", "proto"), filepath.Join(Dir, "app", "desc", "proto"))
	cobra.CheckErr(err)

	// write config.toml
	configTomlFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "config.toml.tpl")))
	cobra.CheckErr(err)

	g, err := genius.NewFromToml(configTomlFile)
	cobra.CheckErr(err)

	switch ConfigType {
	case "toml":
		err = os.WriteFile(filepath.Join(Dir, "config.toml"), configTomlFile, 0o644)
		cobra.CheckErr(err)
	case "yaml":
		yaml, err := g.EncodeToYaml()
		cobra.CheckErr(err)
		err = os.WriteFile(filepath.Join(Dir, "config.yaml"), yaml, 0o644)
		cobra.CheckErr(err)
	case "json":
		json, err := g.EncodeToJSON()
		cobra.CheckErr(err)
		err = os.WriteFile(filepath.Join(Dir, "config.json"), json, 0o644)
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
	err = os.WriteFile(filepath.Join(Dir, "app", "internal", "config", "config.go"), configGoFile, 0o644)
	cobra.CheckErr(err)
	// ################# end gen config ###################

	// ################# start gen middlewares ###################
	// write app/middlewares/response.go
	err = os.MkdirAll(filepath.Join(Dir, "app", "middlewares"), 0o755)
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "middlewares", "response.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "response.go.tpl")), 0o644)
	cobra.CheckErr(err)

	// write app/middlewares/errors.go
	err = os.WriteFile(filepath.Join(Dir, "app", "middlewares", "errors.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "errors.go.tpl")), 0o644)
	cobra.CheckErr(err)

	// write app/middlewares/grpc_rate_limit.go
	err = os.WriteFile(filepath.Join(Dir, "app", "middlewares", "grpc_rate_limit.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "grpc_rate_limit.go.tpl")), 0o644)
	cobra.CheckErr(err)

	// write app/middlewares/logs.go
	logsMiddlewareFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":    APP,
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", "logs.go.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "middlewares", "logs.go"), logsMiddlewareFile, 0o644)
	cobra.CheckErr(err)

	// ################# end gen middlewares ###################

	// write app/internal/handler/myroutes.go
	_ = os.MkdirAll(filepath.Join(Dir, "app", "internal", "handler"), 0o755)
	myroutesFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "handler", "myroutes.go.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "internal", "handler", "myroutes.go"), myroutesFile, 0o644)
	cobra.CheckErr(err)

	// write app/internal/handler/myhandler.go
	myhandlerFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "handler", "myhandler.go.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "app", "internal", "handler", "myhandler.go"), myhandlerFile, 0o644)
	cobra.CheckErr(err)

	// write Dockerfile
	dockerFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":        APP,
		"ConfigType": ConfigType,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "Dockerfile"), dockerFile, 0o644)
	cobra.CheckErr(err)

	// write Dockerfile-arm64
	dockerArm64File, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":        APP,
		"ConfigType": ConfigType,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile-arm64.tpl")))
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "Dockerfile-arm64"), dockerArm64File, 0o644)
	cobra.CheckErr(err)

	// write .goreleaser.yaml
	goreleaserBytes := embeded.ReadTemplateFile(filepath.Join("jzero", "goreleaser.yaml.tpl"))
	newGoreleaserBytes := bytes.ReplaceAll(goreleaserBytes, []byte("{{ .APP }}"), []byte(APP))
	err = os.WriteFile(filepath.Join(Dir, ".goreleaser.yaml"), newGoreleaserBytes, 0o644)
	cobra.CheckErr(err)

	// write Taskfile.yml
	err = os.WriteFile(filepath.Join(Dir, "Taskfile.yml"), embeded.ReadTemplateFile(filepath.Join("jzero", "Taskfile.yml.tpl")), 0o644)
	cobra.CheckErr(err)

	err = embeded.WriteTemplateDir(filepath.Join("go-zero"), filepath.Join(embeded.Home, "go-zero"))
	cobra.CheckErr(err)

	return nil
}
