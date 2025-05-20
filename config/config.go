package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

// C global command flags
var C Config

type Config struct {
	// global flags
	Debug bool `mapstructure:"debug"`

	Hooks HooksConfig `mapstructure:"hooks"`

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

	// migrate command
	Migrate MigrateConfig `mapstructure:"migrate"`

	// format command
	Format FormatConfig `mapstructure:"format"`
}

type NewConfig struct {
	Name               string   `mapstructure:"name"`                 // 新建项目名称
	Home               string   `mapstructure:"home"`                 // 新建项目使用的模板文件目录
	Module             string   `mapstructure:"module"`               // 新建的项目的 go module
	Mono               bool     `mapstructure:"mono"`                 // 是否是 mono 项目(即在一个mod项目之下, 但该项目本身无 go.mod 文件)
	Output             string   `mapstructure:"output"`               // 输出到的目录
	Remote             string   `mapstructure:"remote"`               // 远程仓库地址
	Cache              bool     `mapstructure:"cache"`                // 当使用远程仓库时是否使用缓存
	RemoteAuthUsername string   `mapstructure:"remote-auth-username"` // 远程仓库的认证用户名
	RemoteAuthPassword string   `mapstructure:"remote-auth-password"` // 远程仓库的认证密码
	Frame              string   `mapstructure:"frame"`                // 使用 jzero 内置的框架
	Branch             string   `mapstructure:"branch"`               // 使用远程模板仓库的某个分支
	Local              string   `mapstructure:"local"`                // 使用本地模板与 branch 对应
	Features           []string `mapstructure:"features"`             // 新建项目使用哪些特性, 灵活构建模板
}

type GenConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	// global flags
	Home string `mapstructure:"home"`

	Style string `mapstructure:"style"`

	// code style flags
	ChangeLogicTypes bool `mapstructure:"change-logic-types"`

	RpcStylePatch   bool `mapstructure:"rpc-style-patch"`
	RegenApiHandler bool `mapstructure:"regen-api-handler"`

	// git flags
	GitChange bool `mapstructure:"git-change"`

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
	Crud       GenCrudConfig       `mapstructure:"crud"`
}

type GenSdkConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	Desc         []string `mapstructure:"desc"`
	DescIgnore   []string `mapstructure:"desc-ignore"`
	Scope        string   `mapstructure:"scope"`
	WrapResponse bool     `mapstructure:"wrap-response"`
	Output       string   `mapstructure:"output"`
	Language     string   `mapstructure:"language"`
	GoVersion    string   `mapstructure:"goVersion"`
	GoModule     string   `mapstructure:"goModule"`
	GoPackage    string   `mapstructure:"goPackage"`
	Mono         bool     `mapstructure:"mono"`
}

type GenSwaggerConfig struct {
	Desc       []string `mapstructure:"desc"`
	DescIgnore []string `mapstructure:"desc-ignore"`
	Output     string   `mapstructure:"output"`
	Route2Code bool     `mapstructure:"route2code"`
	Merge      bool     `mapstructure:"merge"`
}

type GenZrpcclientConfig struct {
	Hooks      HooksConfig `mapstructure:"hooks"`
	Desc       []string    `mapstructure:"desc"`
	DescIgnore []string    `mapstructure:"desc-ignore"`
	PbDir      string      `mapstructure:"pb-dir"`
	ClientDir  string      `mapstructure:"client-dir"`
	Scope      string      `mapstructure:"scope"`
	Output     string      `mapstructure:"output"`
	GoVersion  string      `mapstructure:"goVersion"`
	GoModule   string      `mapstructure:"goModule"`
	GoPackage  string      `mapstructure:"goPackage"`
}

type GenDocsConfig struct {
	Desc       []string `mapstructure:"desc"`
	DescIgnore []string `mapstructure:"desc-ignore"`
	Output     string   `mapstructure:"output"`
	Format     string   `mapstructure:"format"`
}

type GenCrudConfig struct {
	// todo: add flag
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

	Delete ServerlessDeleteConfig `mapstructure:"delete"`
}

type ServerlessDeleteConfig struct {
	Plugin []string `mapstructure:"plugin"`
}

type MigrateConfig struct {
	Source   string `mapstructure:"source"`
	Database string `mapstructure:"database"`
}

type FormatConfig struct {
	GitChange   bool `mapstructure:"git-change"`
	DisplayDiff bool `mapstructure:"display-diff"`
}

type HooksConfig struct {
	Before []string `mapstructure:"before"`
	After  []string `mapstructure:"after"`
}

func (c *Config) HomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func (c *Config) Wd() string {
	wd, _ := os.Getwd()
	return wd
}

func (c *Config) ProtoDir() string {
	return filepath.Join("desc", "proto")
}

func (c *Config) ApiDir() string {
	return filepath.Join("desc", "api")
}

func (c *Config) SqlDir() string {
	return filepath.Join("desc", "sql")
}

func (c *Config) SwaggerDir() string {
	return filepath.Join("desc", "swagger")
}

var (
	once         sync.Once
	goctlVersion *version.Version
)

func (c *Config) GoctlVersion() *version.Version {
	once.Do(func() {
		goctlVersionResp, err := execx.Run("goctl -v", "")
		if err != nil {
			panic(err)
		}

		logx.Debugf("goctl version: %s", goctlVersionResp)
		versionInfo := strings.Split(goctlVersionResp, " ")
		if len(versionInfo) >= 3 {
			goctlVersion, err = version.NewVersion(versionInfo[2])
			if err != nil {
				panic(err)
			}
		}
	})

	return goctlVersion
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
