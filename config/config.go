package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// C global command flags
var C Config

type Config struct {
	Syntax string `mapstructure:"syntax"`

	// global flags
	Debug bool `mapstructure:"debug"`

	// new command
	New NewConfig `mapstructure:"new"`
	// gen command
	Gen GenConfig `mapstructure:"gen"`
	// ivm command
	Ivm IvmConfig `mapstructure:"ivm"`
	// template command
	Template TemplateConfig `mapstructure:"template"`
	// upgrade command
	Upgrade UpgradeConfig `mapstructure:"upgrade"`
}

type NewConfig struct {
	Home         string   `mapstructure:"home"`
	Module       string   `mapstructure:"module"`
	Output       string   `mapstructure:"output"`
	Remote       string   `mapstructure:"remote"`
	Cache        bool     `mapstructure:"cache"`
	Branch       string   `mapstructure:"branch"`
	WithTemplate bool     `mapstructure:"with-template"`
	Style        string   `mapstructure:"style"`
	Features     []string `mapstructure:"features"`
}

type GenConfig struct {
	// global flags
	Home  string `mapstructure:"home"`
	Style string `mapstructure:"style"`

	ChangeReplaceTypes        bool     `mapstructure:"change-replace-types"`
	RemoveSuffix              bool     `mapstructure:"remove-suffix"`
	ModelMysqlIgnoreColumns   []string `mapstructure:"model-mysql-ignore-columns"`
	ModelMysqlDatasource      bool     `mapstructure:"model-mysql-datasource"`
	ModelMysqlDatasourceUrl   string   `mapstructure:"model-mysql-datasource-url"`
	ModelMysqlDatasourceTable []string `mapstructure:"model-mysql-datasource-table"`
	ModelMysqlCache           bool     `mapstructure:"model-mysql-cache"`
	ModelMysqlCachePrefix     string   `mapstructure:"model-mysql-cache-prefix"`

	Sdk        GenSdkConfig        `mapstructure:"sdk"`
	Swagger    GenSwaggerConfig    `mapstructure:"swagger"`
	Zrpcclient GenZrpcclientConfig `mapstructure:"zrpcclient"`
}

type GenSdkConfig struct {
	Scope        string `mapstructure:"scope"`
	ApiDir       string `mapstructure:"api-dir"`
	ProtoDir     string `mapstructure:"proto-dir"`
	WrapResponse bool   `mapstructure:"wrap-response"`
	Output       string `mapstructure:"output"`
	Language     string `mapstructure:"language"`
	GoModule     string `mapstructure:"goModule"`
	GoPackage    string `mapstructure:"goPackage"`
}

type GenSwaggerConfig struct {
	Output   string `mapstructure:"output"`
	ApiDir   string `mapstructure:"api-dir"`
	ProtoDir string `mapstructure:"proto-dir"`
}

type GenZrpcclientConfig struct {
	Scope     string `mapstructure:"scope"`
	Output    string `mapstructure:"output"`
	GoModule  string `mapstructure:"goModule"`
	GoPackage string `mapstructure:"goPackage"`
}

type IvmConfig struct {
	// global flags
	Version string `mapstructure:"version"`

	Init IvmInitConfig `mapstructure:"init"`
	Add  IvmAddConfig  `mapstructure:"add"`
}

type IvmInitConfig struct {
	Style              string `mapstructure:"style"`
	RemoveSuffix       bool   `mapstructure:"remove-suffix"`
	ChangeReplaceTypes bool   `mapstructure:"change-replace-types"`
}

type IvmAddConfig struct {
	Api   IvmAddApiConfig   `mapstructure:"api"`
	Proto IvmAddProtoConfig `mapstructure:"proto"`
}

type IvmAddApiConfig struct {
	Name     string
	Group    string
	Handlers []string
}

type IvmAddProtoConfig struct {
	Methods  []string
	Name     string
	Services []string
}

type TemplateConfig struct {
	Init  TemplateInitConfig  `mapstructure:"init"`
	Build TemplateBuildConfig `mapstructure:"build"`
}

type TemplateInitConfig struct {
	Output string `mapstructure:"output"`
	Remote string `mapstructure:"remote"`
	Branch string `mapstructure:"branch"`
}

type TemplateBuildConfig struct {
	Output     string `mapstructure:"output"`
	WorkingDir string `mapstructure:"working-dir"`
	Name       string `mapstructure:"name"`
}

type UpgradeConfig struct {
	Channel string `mapstructure:"channel"`
}

func (gc *GenSwaggerConfig) Wd() string {
	wd, _ := os.Getwd()
	return wd
}

func (gc *GenConfig) Wd() string {
	wd, _ := os.Getwd()
	return wd
}

func SetConfig(command string, flagSet *pflag.FlagSet) error {
	flagSet.VisitAll(func(flag *pflag.Flag) {
		if command == "" {
			err := viper.BindPFlag(flag.Name, flag)
			if err != nil {
				panic(err)
			}
		} else {
			err := viper.BindPFlag(fmt.Sprintf("%s.%s", command, flag.Name), flag)
			if err != nil {
				panic(err)
			}
		}
	})

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
