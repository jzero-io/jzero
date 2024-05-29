package new

import (
	"path/filepath"

	"github.com/jaronnie/genius"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroConfig struct {
	TemplateData map[string]interface{}
}

func (jc *JzeroConfig) New() error {
	configYamlFile, err := templatex.ParseTemplate(jc.TemplateData, embeded.ReadTemplateFile(filepath.Join("jzero", "config.yaml.tpl")))
	if err != nil {
		return err
	}

	g, err := genius.NewFromYaml(configYamlFile)
	if err != nil {
		return err
	}

	switch ConfigType {
	case "toml":
		configTomlFile, err := g.EncodeToToml()
		if err != nil {
			return err
		}
		err = checkWrite(filepath.Join(Output, "config.toml"), configTomlFile)
		if err != nil {
			return err
		}
	case "yaml":
		err = checkWrite(filepath.Join(Output, "config.yaml"), configYamlFile)
		if err != nil {
			return err
		}
	case "json":
		configJsonFile, err := g.EncodeToJSON()
		if err != nil {
			return err
		}
		err = checkWrite(filepath.Join(Output, "config.json"), configJsonFile)
		if err != nil {
			return err
		}
	}
	return nil
}
