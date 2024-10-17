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
	Syntax  string `mapstructure:"syntax"`
	Version string `mapstructure:"version"`

	/*
		===============================command flags start========================================
	*/
	// global flags
	Debug bool `mapstructure:"debug"`

	DebugSleepTime int `mapstructure:"debug-sleep-time"`

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

	/*
		==============================command flags end=========================================
	*/
}

type NewConfig struct {
	Home         string   `mapstructure:"home"`
	Module       string   `mapstructure:"module"`
	SubModule    bool     `mapstructure:"submodule"`
	Output       string   `mapstructure:"output"`
	Remote       string   `mapstructure:"remote"`
	Cache        bool     `mapstructure:"cache"`
	Frame        string   `mapstructure:"frame"`
	Branch       string   `mapstructure:"branch"`
	WithTemplate bool     `mapstructure:"with-template"`
	Style        string   `mapstructure:"style"`
	Features     []string `mapstructure:"features"`
}

type GenConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	// global flags

	Home             string `mapstructure:"home"`
	Style            string `mapstructure:"style"`
	Branch           string `mapstructure:"branch"`
	SplitApiTypesDir bool   `mapstructure:"split-api-types-dir"`

	// code style flags

	ChangeLogicTypes bool `mapstructure:"change-logic-types"`
	RpcStylePatch    bool `mapstructure:"rpc-style-patch"`
	RemoveSuffix     bool `mapstructure:"remove-suffix"`
	RegenApiHandler  bool `mapstructure:"regen-api-handler"`

	// git flags

	GitChange          bool   `mapstructure:"git-change"`
	ApiGitChangePath   string `mapstructure:"api-git-change-path"`
	ModelGitChangePath string `mapstructure:"model-git-change-path"`
	ProtoGitChangePath string `mapstructure:"proto-git-change-path"`

	// model flags

	ModelMysqlStrict          bool     `mapstructure:"model-mysql-strict"`
	ModelMysqlIgnoreColumns   []string `mapstructure:"model-mysql-ignore-columns"`
	ModelMysqlDatasource      bool     `mapstructure:"model-mysql-datasource"`
	ModelMysqlDatasourceUrl   string   `mapstructure:"model-mysql-datasource-url"`
	ModelMysqlDatasourceTable []string `mapstructure:"model-mysql-datasource-table"`
	ModelMysqlCache           bool     `mapstructure:"model-mysql-cache"`
	ModelMysqlCachePrefix     string   `mapstructure:"model-mysql-cache-prefix"`

	// gen code flags
	Desc []string `mapstructure:"desc"`

	Sdk        GenSdkConfig        `mapstructure:"sdk"`
	Swagger    GenSwaggerConfig    `mapstructure:"swagger"`
	Zrpcclient GenZrpcclientConfig `mapstructure:"zrpcclient"`
	Docs       GenDocsConfig       `mapstructure:"docs"`
}

type GenSdkConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	Scope        string `mapstructure:"scope"`
	ApiDir       string `mapstructure:"api-dir"`
	ProtoDir     string `mapstructure:"proto-dir"`
	WrapResponse bool   `mapstructure:"wrap-response"`
	Output       string `mapstructure:"output"`
	Language     string `mapstructure:"language"`
	GoVersion    string `mapstructure:"goVersion"`
	GoModule     string `mapstructure:"goModule"`
	GoPackage    string `mapstructure:"goPackage"`
}

type GenSwaggerConfig struct {
	Output   string `mapstructure:"output"`
	ApiDir   string `mapstructure:"api-dir"`
	ProtoDir string `mapstructure:"proto-dir"`
}

type GenZrpcclientConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	Scope     string `mapstructure:"scope"`
	Output    string `mapstructure:"output"`
	GoVersion string `mapstructure:"goVersion"`
	GoModule  string `mapstructure:"goModule"`
	GoPackage string `mapstructure:"goPackage"`
}

type GenDocsConfig struct {
	Output   string `mapstructure:"output"`
	Format   string `mapstructure:"format"`
	ApiDir   string `mapstructure:"api-dir"`
	ProtoDir string `mapstructure:"proto-dir"`
}

type IvmConfig struct {
	// global flags
	Version string `mapstructure:"version"`

	SplitApiTypesDir bool `mapstructure:"split-api-types-dir"`

	Init IvmInitConfig `mapstructure:"init"`
	Add  IvmAddConfig  `mapstructure:"add"`
}

type IvmInitConfig struct {
	Style            string `mapstructure:"style"`
	RemoveSuffix     bool   `mapstructure:"remove-suffix"`
	ChangeLogicTypes bool   `mapstructure:"change-logic-types"`
}

type IvmAddConfig struct {
	Api   IvmAddApiConfig   `mapstructure:"api"`
	Proto IvmAddProtoConfig `mapstructure:"proto"`
}

type IvmAddApiConfig struct {
	Name     string   `mapstructure:"name"`
	Group    string   `mapstructure:"group"`
	Handlers []string `mapstructure:"handlers"`
}

type IvmAddProtoConfig struct {
	Methods  []string `mapstructure:"methods"`
	Name     string   `mapstructure:"name"`
	Services []string `mapstructure:"services"`
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

type HooksConfig struct {
	Before []string `mapstructure:"before"`
	After  []string `mapstructure:"after"`
}

func (c *Config) Wd() string {
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
