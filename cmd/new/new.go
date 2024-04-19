package new

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"

	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"
)

var (
	Module string
	Dir    string
	APP    string
)

func NewProject(_ *cobra.Command, _ []string) error {
	// mkdir output
	err := os.MkdirAll(Dir, 0o755)
	cobra.CheckErr(err)
	// go mod init
	_, err = execx.Run(fmt.Sprintf("go mod init %s", Module), Dir)
	cobra.CheckErr(err)
	// touch main.go
	mainFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "main.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "main.go"), mainFile, 0o644)
	cobra.CheckErr(err)
	// mkdir cmd dir
	err = os.MkdirAll(filepath.Join(Dir, "cmd"), 0o755)
	cobra.CheckErr(err)
	// touch cmd/root.go
	rootCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "root.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "cmd", "root.go"), rootCmdFile, 0o644)
	cobra.CheckErr(err)
	// touch cmd/daemon.go
	daemonCmdFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", "daemon.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "cmd", "daemon.go"), daemonCmdFile, 0o644)
	cobra.CheckErr(err)
	// mkdir daemon dir
	err = os.MkdirAll(filepath.Join(Dir, "daemon"), 0o755)
	cobra.CheckErr(err)
	// touch daemon/daemon.go
	daemonFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "daemon.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "daemon.go"), daemonFile, 0o644)
	cobra.CheckErr(err)

	// touch daemon/zrpc.go
	zrpcFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
		"APP":    APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "zrpc.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "zrpc.go"), zrpcFile, 0o644)
	cobra.CheckErr(err)

	// mkdir api, proto dir
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "desc", "proto"), 0o755)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "desc", "api"), 0o755)
	cobra.CheckErr(err)
	// touch daemon/desc/api/{{.APP}}.api
	err = os.WriteFile(filepath.Join(Dir, "daemon", "desc", "api", APP+".api"), embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "desc", "api", "jzero.api.tpl")), 0o644)
	cobra.CheckErr(err)
	// touch daemon/desc/api/hello.api
	helloApiFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "desc", "api", "hello.api.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "desc", "api", "hello.api"), helloApiFile, 0o644)
	cobra.CheckErr(err)

	// write proto dir
	err = embeded.WriteTemplateDir(filepath.Join("jzero", "daemon", "desc", "proto"), filepath.Join(Dir, "daemon", "desc", "proto"))
	cobra.CheckErr(err)

	// write config.toml
	configTomlFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "config.toml.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "config.toml"), configTomlFile, 0o644)
	cobra.CheckErr(err)

	// ################# start gen config ###################
	// write daemon/internal/config/config.go
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "internal", "config"), 0o755)
	cobra.CheckErr(err)

	configGoFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "internal", "config", "config.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "internal", "config", "config.go"), configGoFile, 0o644)
	cobra.CheckErr(err)
	// ################# end gen config ###################

	// ################# start gen middlewares ###################
	// write daemon/middlewares/response.go
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "middlewares"), 0o755)
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "daemon", "middlewares", "response.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "middlewares", "response.go.tpl")), 0o644)
	cobra.CheckErr(err)

	// write daemon/middlewares/errors.go
	err = os.WriteFile(filepath.Join(Dir, "daemon", "middlewares", "errors.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "middlewares", "errors.go.tpl")), 0o644)
	cobra.CheckErr(err)

	// ################# end gen middlewares ###################

	// write daemon/internal/handler/myroutes.go
	_ = os.MkdirAll(filepath.Join(Dir, "daemon", "internal", "handler"), 0o755)
	myroutesFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "internal", "handler", "myroutes.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "internal", "handler", "myroutes.go"), myroutesFile, 0o644)
	cobra.CheckErr(err)

	// write daemon/internal/handler/myhandler.go
	myhandlerFile, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "internal", "handler", "myhandler.go.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "internal", "handler", "myhandler.go"), myhandlerFile, 0o644)
	cobra.CheckErr(err)

	// write Dockerfile
	dockerFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "Dockerfile"), dockerFile, 0o644)
	cobra.CheckErr(err)

	// write Dockerfile-arm64
	dockerArm64File, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile-arm64.tpl")))
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

	return nil
}
