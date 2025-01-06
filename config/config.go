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

	// serverless command
	Serverless ServerlessConfig `mapstructure:"serverless"`

	/*
		==============================command flags end=========================================
	*/
}

type NewConfig struct {
	Core     bool     `mapstructure:"core"`
	Home     string   `mapstructure:"home"`     // 新建项目使用的模板文件目录
	Module   string   `mapstructure:"module"`   // 新建的项目的 go module
	Mono     bool     `mapstructure:"mono"`     // 是否是 mono 项目(即在一个mod项目之下, 但该项目本身无 go.mod 文件)
	Output   string   `mapstructure:"output"`   // 输出到的目录
	Remote   string   `mapstructure:"remote"`   // 远程仓库地址
	Frame    string   `mapstructure:"frame"`    // 使用 jzero 内置的框架
	Branch   string   `mapstructure:"branch"`   // 使用远程模板仓库的某个分支
	Local    string   `mapstructure:"local"`    // 使用本地模板与 branch 对应
	Features []string `mapstructure:"features"` // 新建项目使用哪些特性, 灵活构建模板
}

type GenConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	// global flags
	Home string `mapstructure:"home"`

	Style            string `mapstructure:"style"`
	SplitApiTypesDir bool   `mapstructure:"split-api-types-dir"`

	// code style flags
	ChangeLogicTypes bool `mapstructure:"change-logic-types"`

	RpcStylePatch   bool `mapstructure:"rpc-style-patch"`
	RegenApiHandler bool `mapstructure:"regen-api-handler"`

	// git flags
	GitChange bool `mapstructure:"git-change"`

	ApiGitChangePath   string `mapstructure:"api-git-change-path"`
	ModelGitChangePath string `mapstructure:"model-git-change-path"`
	ProtoGitChangePath string `mapstructure:"proto-git-change-path"`

	// model flags
	ModelMysqlStrict bool `mapstructure:"model-mysql-strict"`

	ModelMysqlIgnoreColumns   []string `mapstructure:"model-mysql-ignore-columns"`
	ModelMysqlDDLDatabase     string   `mapstructure:"model-mysql-ddl-database"`
	ModelMysqlDatasource      bool     `mapstructure:"model-mysql-datasource"`
	ModelMysqlDatasourceUrl   string   `mapstructure:"model-mysql-datasource-url"`
	ModelMysqlDatasourceTable []string `mapstructure:"model-mysql-datasource-table"`
	ModelMysqlCache           bool     `mapstructure:"model-mysql-cache"`
	ModelMysqlCachePrefix     string   `mapstructure:"model-mysql-cache-prefix"`
	GenMysqlCreateTableDDL    bool     `mapstructure:"gen-mysql-create-table-ddl"`

	// rpc flags
	RpcClient bool `mapstructure:"rpc-client"`

	// gen code flags
	Desc []string `mapstructure:"desc"`

	DescIgnore []string `mapstructure:"desc-ignore"`

	// other
	Route2Code bool

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
	Mono         bool   `mapstructure:"mono"`
}

type GenSwaggerConfig struct {
	Output   string `mapstructure:"output"`
	ApiDir   string `mapstructure:"api-dir"`
	ProtoDir string `mapstructure:"proto-dir"`
}

type GenZrpcclientConfig struct {
	Hooks     HooksConfig `mapstructure:"hooks"`
	PbDir     string      `mapstructure:"pb-dir"`
	ClientDir string      `mapstructure:"client-dir"`
	Scope     string      `mapstructure:"scope"`
	Output    string      `mapstructure:"output"`
	GoVersion string      `mapstructure:"goVersion"`
	GoModule  string      `mapstructure:"goModule"`
	GoPackage string      `mapstructure:"goPackage"`
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

	Init IvmInitConfig `mapstructure:"init"`
	Add  IvmAddConfig  `mapstructure:"add"`
}

type IvmInitConfig struct {
	Style            string `mapstructure:"style"`
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
	Output     string   `mapstructure:"output"`
	WorkingDir string   `mapstructure:"working-dir"`
	Name       string   `mapstructure:"name"`
	Ignore     []string `mapstructure:"ignore"`
}

type UpgradeConfig struct {
	Channel string `mapstructure:"channel"`
}

type ServerlessConfig struct {
	Home string `mapstructure:"home"` // 使用的模板文件目录

	New    NewConfig              `mapstructure:"new"`
	Delete ServerlessDeleteConfig `mapstructure:"delete"`
}

type ServerlessNewConfig struct {
	Module string `mapstructure:"module"` // 新建的项目的 go module
	Remote string `mapstructure:"remote"` // 远程仓库地址
	Frame  string `mapstructure:"frame"`  // 使用 jzero 内置的框架
	Branch string `mapstructure:"branch"` // 使用远程模板仓库的某个分支
	Local  string `mapstructure:"local"`  // 使用本地模板与 branch 对应
}

type ServerlessDeleteConfig struct {
	Plugin []string `mapstructure:"plugin"`
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

	viper.SetEnvPrefix("JZERO")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
