package new

import (
	"bytes"
	"fmt"
	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"os"
	"path/filepath"
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
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "proto"), 0o755)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "api"), 0o755)
	cobra.CheckErr(err)
	// touch daemon/api/{{.APP}}.api
	err = os.WriteFile(filepath.Join(Dir, "daemon", "api", APP+".api"), embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "api", "jzero.api.tpl")), 0o644)
	cobra.CheckErr(err)
	// touch daemon/api/hello.api
	helloApiFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "api", "hello.api.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "daemon", "api", "hello.api"), helloApiFile, 0o644)
	cobra.CheckErr(err)

	// write proto dir
	err = embeded.WriteTemplateDir(filepath.Join("jzero", "daemon", "proto"), filepath.Join(Dir, "daemon", "proto"))
	cobra.CheckErr(err)

	// write config.toml
	configFile, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": APP,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "config.toml.tpl")))
	err = os.WriteFile(filepath.Join(Dir, "config.toml"), configFile, 0o644)
	cobra.CheckErr(err)

	// write .template
	err = embeded.WriteTemplateDir("go-zero", filepath.Join(Dir, ".template", "go-zero"))
	cobra.CheckErr(err)
	// replace .template/go-zero/api/handler.tpl
	// 暂时特殊处理: https://github.com/zeromicro/go-zero/pull/4071
	newHandlerBytes := bytes.ReplaceAll(embeded.ReadTemplateFile(filepath.Join(filepath.Join("go-zero", "api", "handler.tpl"))), []byte("github.com/jaronnie/jzero"), []byte(Module))
	err = os.WriteFile(filepath.Join(Dir, ".template", "go-zero", "api", "handler.tpl"), newHandlerBytes, 0o644)
	cobra.CheckErr(err)

	// write daemon/pkg/response/response.go
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "pkg", "response"), 0o755)
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "daemon", "pkg", "response", "response.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "pkg", "response", "response.go.tpl")), 0o644)
	cobra.CheckErr(err)

	// write daemon/internal/config/config.go
	err = os.MkdirAll(filepath.Join(Dir, "daemon", "internal", "config"), 0o755)
	cobra.CheckErr(err)
	err = os.WriteFile(filepath.Join(Dir, "daemon", "internal", "config", "config.go"), embeded.ReadTemplateFile(filepath.Join("jzero", "daemon", "internal", "config", "config.go.tpl")), 0o644)
	cobra.CheckErr(err)

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
