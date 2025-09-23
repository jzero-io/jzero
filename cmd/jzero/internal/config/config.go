package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/a8m/envsubst"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/protoc"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/protocgengo"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/protocgengogrpc"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"gopkg.in/yaml.v3"

	"github.com/jzero-io/jzero/cmd/jzero/internal/hooks"
)

// C global command flags
var C Config

var (
	CfgFile    string
	CfgEnvFile string
)

type Config struct {
	// global flags
	Debug bool `mapstructure:"debug"`

	// Register tpl val
	RegisterTplVal []string `mapstructure:"register-tpl-val"`

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
	Name                 string   `mapstructure:"name"`                  // 新建项目名称
	Home                 string   `mapstructure:"home"`                  // 新建项目使用的模板文件目录
	Module               string   `mapstructure:"module"`                // 新建的项目的 go module
	Mono                 bool     `mapstructure:"mono"`                  // 是否是 mono 项目(即在一个mod项目之下, 但该项目本身无 go.mod 文件)
	Serverless           bool     `mapstructure:"serverless"`            // 是否是 serverless 插件
	Output               string   `mapstructure:"output"`                // 输出到的目录
	Remote               string   `mapstructure:"remote"`                // 远程仓库地址
	RemoteTimeout        int      `mapstructure:"remote-timeout"`        // 远程仓库超时时间, 单位秒
	Cache                bool     `mapstructure:"cache"`                 // 当使用远程仓库时是否使用缓存
	Gen                  bool     `mapstructure:"gen"`                   // 新建项目后是否自动执行 gen 命令
	RemoteAuthUsername   string   `mapstructure:"remote-auth-username"`  // 远程仓库的认证用户名
	RemoteAuthPassword   string   `mapstructure:"remote-auth-password"`  // 远程仓库的认证密码
	Frame                string   `mapstructure:"frame"`                 // 使用 jzero 内置的框架
	Style                string   `mapstructure:"style"`                 // 代码风格
	Branch               string   `mapstructure:"branch"`                // 使用远程模板仓库的某个分支
	Local                string   `mapstructure:"local"`                 // 使用本地模板与 branch 对应
	Features             []string `mapstructure:"features"`              // 新建项目使用哪些特性, 灵活构建模板
	Ignore               []string `mapstructure:"ignore"`                // 忽略哪些文件或目录
	IgnoreExtra          []string `mapstructure:"ignore-extra"`          // 忽略哪些额外的文件或目录
	ExecutableExtensions []string `mapstructure:"executable-extensions"` // 可执行文件的后缀
}

type GenConfig struct {
	// Hooks
	Hooks HooksConfig `mapstructure:"hooks"`

	// gen global flags
	Home string `mapstructure:"home"`

	Style      string   `mapstructure:"style"`
	Desc       []string `mapstructure:"desc"`
	DescIgnore []string `mapstructure:"desc-ignore"`

	// gen self flags
	GitChange bool `mapstructure:"git-change"`

	Route2Code   bool
	RpcClient    bool     `mapstructure:"rpc-client"`
	ProtoInclude []string `mapstructure:"proto-include"`

	// model flag
	ModelDriver string `mapstructure:"model-driver"`

	ModelStrict          bool     `mapstructure:"model-strict"`
	ModelIgnoreColumns   []string `mapstructure:"model-ignore-columns"`
	ModelSchema          string   `mapstructure:"model-schema"`
	ModelDatasource      bool     `mapstructure:"model-datasource"`
	ModelDatasourceUrl   []string `mapstructure:"model-datasource-url"`
	ModelDatasourceTable []string `mapstructure:"model-datasource-table"`
	ModelCache           bool     `mapstructure:"model-cache"`
	ModelCacheTable      []string `mapstructure:"model-cache-table"`
	ModelCachePrefix     string   `mapstructure:"model-cache-prefix"`
	ModelCreateTableDDL  bool     `mapstructure:"model-create-table-ddl"`

	// mongo flags
	MongoType []string `mapstructure:"mongo-type"`

	MongoCache       bool     `mapstructure:"mongo-cache"`
	MongoCachePrefix string   `mapstructure:"mongo-cache-prefix"`
	MongoCacheType   []string `mapstructure:"mongo-cache-type"`

	// Sub command
	Sdk GenSdkConfig `mapstructure:"sdk"`

	Swagger    GenSwaggerConfig    `mapstructure:"swagger"`
	Zrpcclient GenZrpcclientConfig `mapstructure:"zrpcclient"`
	Docs       GenDocsConfig       `mapstructure:"docs"`
}

type GenSdkConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	Desc       []string `mapstructure:"desc"`
	DescIgnore []string `mapstructure:"desc-ignore"`
	Output     string   `mapstructure:"output"`
	Language   string   `mapstructure:"language"`
	GoVersion  string   `mapstructure:"goVersion"`
	GoModule   string   `mapstructure:"goModule"`
	GoPackage  string   `mapstructure:"goPackage"`
	Mono       bool     `mapstructure:"mono"`
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
	Output     string      `mapstructure:"output"`
	GoVersion  string      `mapstructure:"goVersion"`
	GoModule   string      `mapstructure:"goModule"`
	GoPackage  string      `mapstructure:"goPackage"`
	Mono       bool        `mapstructure:"mono"`
}

type GenDocsConfig struct {
	Desc       []string `mapstructure:"desc"`
	DescIgnore []string `mapstructure:"desc-ignore"`
	Output     string   `mapstructure:"output"`
	Format     string   `mapstructure:"format"`
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
	goctlVersionOnce       sync.Once
	protocVersionOnce      sync.Once
	protocGenGoVersionOnce sync.Once
	protocGenGoGrpcOnce    sync.Once
	protocGenOpenapiv2     sync.Once
	protocGenDocOnce       sync.Once

	protocVersion             *version.Version
	goctlVersion              *version.Version
	protocGenGoVersion        *version.Version
	protocGenGoGrpcVersion    *version.Version
	protocGenOpenapiv2Version *version.Version
	protocGenDocVersion       *version.Version
)

func (c *Config) ProtocVersion() *version.Version {
	protocVersionOnce.Do(func() {
		protocVersionResp, err := protoc.Version()
		if err != nil {
			logx.Debugf("get protoc version failed, %v", err)
			return
		}
		logx.Debugf("protoc version: %s", protocVersionResp)
		protocVersion, err = version.NewVersion(strings.TrimPrefix(protocVersionResp, "libprotoc "))
		if err != nil {
			logx.Debugf("new protoc version failed, %v", err)
			return
		}
	})
	return protocVersion
}

func (c *Config) GoctlVersion() *version.Version {
	goctlVersionOnce.Do(func() {
		goctlVersionResp, err := execx.Run("goctl -v", "")
		if err != nil {
			logx.Debugf("get goctl version failed, %v", err)
			return
		}

		logx.Debugf("goctl version: %s", goctlVersionResp)
		versionInfo := strings.Split(goctlVersionResp, " ")
		if len(versionInfo) >= 3 {
			goctlVersion, err = version.NewVersion(strings.TrimPrefix(versionInfo[2], "v"))
			if err != nil {
				logx.Debugf("new goctl version failed, %v", err)
				return
			}
		}
	})

	return goctlVersion
}

func (c *Config) ProtocGenGoVersion() *version.Version {
	protocGenGoVersionOnce.Do(func() {
		protocGenGoVersionResp, err := protocgengo.Version()
		if err != nil {
			logx.Debugf("get protoc-gen-go version failed, %v", err)
			return
		}
		logx.Debugf("protoc-gen-go version: %s", protocGenGoVersionResp)
		protocGenGoVersion, err = version.NewVersion(strings.TrimPrefix(protocGenGoVersionResp, "v"))
		if err != nil {
			logx.Debugf("new protoc-gen-go version failed, %v", err)
			return
		}
	})

	return protocGenGoVersion
}

func (c *Config) ProtocGenGoGrpcVersion() *version.Version {
	protocGenGoGrpcOnce.Do(func() {
		protocGenGoGrpcVersionResp, err := protocgengogrpc.Version()
		if err != nil {
			logx.Debugf("get protoc-gen-go-grpc version failed, %v", err)
			return
		}
		logx.Debugf("protoc-gen-go-grpc version: %s", protocGenGoGrpcVersionResp)
		protocGenGoGrpcVersion, err = version.NewVersion(strings.TrimPrefix(protocGenGoGrpcVersionResp, "v"))
		if err != nil {
			logx.Debugf("new protoc-gen-go-grpc version failed, %v", err)
			return
		}
	})

	return protocGenGoGrpcVersion
}

func (c *Config) ProtocGenOpenapiv2Version() *version.Version {
	protocGenOpenapiv2.Do(func() {
		protocGenOpenapiv2VersionResp, err := execx.Run("protoc-gen-openapiv2 --version", "")
		if err != nil {
			logx.Debugf("get protoc-gen-openapiv2 version failed, %v", err)
			return
		}
		logx.Debugf("protoc-gen-openapiv2 version: %s", protocGenOpenapiv2VersionResp)
		versionInfo := strings.Split(protocGenOpenapiv2VersionResp, " ")
		if len(versionInfo) >= 2 {
			protocGenOpenapiv2Version, err = version.NewVersion(strings.TrimSuffix(strings.TrimPrefix(versionInfo[1], "v"), ","))
			if err != nil {
				logx.Debugf("new protoc-gen-openapiv2 version failed, %v", err)
				return
			}
		}
	})

	return protocGenOpenapiv2Version
}

func (c *Config) ProtocGenDocVersion() *version.Version {
	protocGenDocOnce.Do(func() {
		protocGenDocVersionResp, err := execx.Run("protoc-gen-doc --version", "")
		if err != nil {
			logx.Debugf("get protoc-gen-doc version failed, %v", err)
			return
		}
		logx.Debugf("protoc-gen-doc version: %s", protocGenDocVersionResp)
		versionInfo := strings.Split(protocGenDocVersionResp, " ")
		if len(versionInfo) >= 3 {
			protocGenDocVersion, err = version.NewVersion(strings.TrimPrefix(versionInfo[2], "v"))
			if err != nil {
				logx.Debugf("new protoc-gen-doc version failed, %v", err)
				return
			}
		}
	})

	return protocGenDocVersion
}

func TraverseCommands(prefix string, cmd *cobra.Command) error {
	err := SetConfig(prefix, cmd.Flags())
	if err != nil {
		return err
	}

	for _, subCommand := range cmd.Commands() {
		newPrefix := fmt.Sprintf("%s.%s", prefix, subCommand.Use)
		if prefix == "" {
			newPrefix = subCommand.Use
		}

		beforeHooks := viper.Get(fmt.Sprintf("%s.hooks.before", newPrefix))
		afterHooks := viper.Get(fmt.Sprintf("%s.hooks.after", newPrefix))

		subCommand.PreRunE = func(cmd *cobra.Command, args []string) error {
			if beforeHooks != nil {
				if h, ok := beforeHooks.(string); ok {
					return hooks.Run(cmd, "Before", newPrefix, strings.Split(h, ","))
				}
				if _, ok := beforeHooks.([]any); ok {
					return hooks.Run(cmd, "Before", newPrefix, viper.GetStringSlice(fmt.Sprintf("%s.hooks.before", newPrefix)))
				}
			}
			return nil
		}
		subCommand.PostRunE = func(cmd *cobra.Command, args []string) error {
			if afterHooks != nil {
				if h, ok := afterHooks.(string); ok {
					return hooks.Run(cmd, "After", newPrefix, strings.Split(h, ","))
				}
				if _, ok := afterHooks.([]any); ok {
					return hooks.Run(cmd, "After", newPrefix, viper.GetStringSlice(fmt.Sprintf("%s.hooks.after", newPrefix)))
				}
			}
			return nil
		}

		err = TraverseCommands(newPrefix, subCommand)
		if err != nil {
			return err
		}
	}

	return nil
}

func ResetConfig() {
	C = Config{}
	unsetEnvVarsWithPrefix("JZERO")
}

func InitConfig(rootCmd *cobra.Command) error {
	if pathx.FileExists(CfgFile) {
		viper.SetConfigFile(CfgFile)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if pathx.FileExists(CfgEnvFile) {
		data, err := envsubst.ReadFile(CfgEnvFile)
		if err != nil {
			return err
		}
		var envMap map[string]any
		err = yaml.Unmarshal(data, &envMap)
		if err != nil {
			return err
		}

		for k, v := range envMap {
			if vs, ok := v.([]any); ok {
				var envs []string
				for _, e := range vs {
					envs = append(envs, cast.ToString(e))
				}
				_ = os.Setenv(k, strings.Join(envs, ","))
			} else {
				_ = os.Setenv(k, cast.ToString(v))
			}
		}
	}

	if err := TraverseCommands("", rootCmd); err != nil {
		return err
	}
	return nil
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

func unsetEnvVarsWithPrefix(prefix string) {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) < 1 {
			continue
		}
		key := pair[0]

		if key == "JZERO_HOOK_TRIGGERED" {
			continue
		}

		if strings.HasPrefix(key, prefix) {
			_ = os.Unsetenv(key)
		}
	}
}
