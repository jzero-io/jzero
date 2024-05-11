package new

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaronnie/genius"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/jzero-io/jzero/embeded"
)

var (
	Module string
	Dir    string
	APP    string
	// ConfigType config type
	ConfigType string
	// Remote templates repo
	Remote string
	Branch string

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

	templateData := map[string]interface{}{
		"Module":     Module,
		"APP":        APP,
		"ConfigType": ConfigType,
	}

	// touch main.go
	mainFile, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "main.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "main.go"), mainFile)
	cobra.CheckErr(err)

	// write cmd dir
	cmdDir := embeded.ReadTemplateDir(filepath.Join("jzero", "cmd"))
	for _, file := range cmdDir {
		cmdFileBytes, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "cmd", file.Name())))
		cobra.CheckErr(err)
		cmdFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Dir, "cmd", cmdFileName), cmdFileBytes)
		cobra.CheckErr(err)
	}

	// write app/*.go
	appDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app"))
	for _, file := range appDir {
		if file.IsDir() {
			continue
		}
		appFileBytes, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", file.Name())))
		cobra.CheckErr(err)
		appFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Dir, "app", appFileName), appFileBytes)
		cobra.CheckErr(err)
	}

	// write proto dir
	protoDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app", "desc", "proto"))
	for _, file := range protoDir {
		if file.IsDir() {
			continue
		}
		protoFileBytes, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "proto", file.Name())))
		cobra.CheckErr(err)
		protoFileName := file.Name()
		err = checkWrite(filepath.Join(Dir, "app", "desc", "proto", protoFileName), protoFileBytes)
		cobra.CheckErr(err)

		if len(protoFileBytes) > 0 {
			if !pathx.FileExists(filepath.Join(Dir, "app", "desc", "proto", "google")) {
				err = embeded.WriteTemplateDir(filepath.Join("jzero", "app", "desc", "proto", "google"), filepath.Join(Dir, "app", "desc", "proto", "google"))
				cobra.CheckErr(err)
			}
		}
	}

	// touch app/desc/api/{{.APP}}.api
	err = checkWrite(filepath.Join(Dir, "app", "desc", "api", APP+".api"), embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", "jzero.api.tpl")))
	cobra.CheckErr(err)
	// touch app/desc/api/hello.api
	helloApiFile, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "desc", "api", "hello.api.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "desc", "api", "hello.api"), helloApiFile)
	cobra.CheckErr(err)

	// write config.yaml
	configYamlFile, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "config.yaml.tpl")))
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

	// write app/internal/config/config.go
	configGoFile, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "internal", "config", "config.go.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "app", "internal", "config", "config.go"), configGoFile)
	cobra.CheckErr(err)

	middlewareDir := embeded.ReadTemplateDir(filepath.Join("jzero", "app", "middlewares"))
	for _, file := range middlewareDir {
		middlewareFileBytes, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "app", "middlewares", file.Name())))
		cobra.CheckErr(err)
		middlewareFileName := strings.TrimRight(file.Name(), ".tpl")
		err = checkWrite(filepath.Join(Dir, "app", "middlewares", middlewareFileName), middlewareFileBytes)
		cobra.CheckErr(err)
	}

	// write Dockerfile
	dockerFile, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, "Dockerfile"), dockerFile)
	cobra.CheckErr(err)

	// write Dockerfile-arm64
	dockerArm64File, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "Dockerfile-arm64.tpl")))
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
	gitignoreFile, err := templatex.ParseTemplate(templateData, embeded.ReadTemplateFile(filepath.Join("jzero", "gitignore.tpl")))
	cobra.CheckErr(err)
	err = checkWrite(filepath.Join(Dir, ".gitignore"), gitignoreFile)
	cobra.CheckErr(err)

	return nil
}

func checkWrite(path string, bytes []byte) error {
	var err error
	if len(bytes) == 0 {
		return nil
	}
	if !pathx.FileExists(filepath.Join(path)) {
		err = os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			return err
		}
	}

	bytesFormat := bytes
	// if is go file. format it
	if filepath.Ext(path) == ".go" {
		bytesFormat, err = gosimports.Process("", bytes, nil)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(path, bytesFormat, 0o644)
}
