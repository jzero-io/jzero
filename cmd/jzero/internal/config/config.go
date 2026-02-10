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

type HooksConfig struct {
	Before []string `mapstructure:"before"`
	After  []string `mapstructure:"after"`
}

type Config struct {
	// config file only
	Hooks HooksConfig `mapstructure:"hooks"`

	// root persistent flags
	Debug bool `mapstructure:"debug"`

	Quiet          bool     `mapstructure:"quiet"`
	DebugSleepTime int      `mapstructure:"debug-sleep-time"`
	RegisterTplVal []string `mapstructure:"register-tpl-val"`
	Home           string   `mapstructure:"home"`
	Style          string   `mapstructure:"style"`

	// new command
	New NewConfig `mapstructure:"new"`

	// add command
	Add AddConfig `mapstructure:"add"`

	// gen command
	Gen GenConfig `mapstructure:"gen"`

	// skills command
	Skills SkillsConfig `mapstructure:"skills"`

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
	Branch               string   `mapstructure:"branch"`                // 使用远程模板仓库的某个分支
	Local                string   `mapstructure:"local"`                 // 使用本地模板与 branch 对应
	Features             []string `mapstructure:"features"`              // 新建项目使用哪些特性, 灵活构建模板
	Ignore               []string `mapstructure:"ignore"`                // 忽略哪些文件或目录
	IgnoreExtra          []string `mapstructure:"ignore-extra"`          // 忽略哪些额外的文件或目录
	ExecutableExtensions []string `mapstructure:"executable-extensions"` // 可执行文件的后缀
}

type GenConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	// gen persistent flags
	// Style: code style
	// Deprecated
	Style string `mapstructure:"style"`

	Desc                    []string `mapstructure:"desc"`
	DescIgnore              []string `mapstructure:"desc-ignore"`
	GitChange               bool     `mapstructure:"git-change"`
	Route2Code              bool
	ProtoInclude            []string `mapstructure:"proto-include"`
	RpcClient               bool     `mapstructure:"rpc-client"`
	ModelDriver             string   `mapstructure:"model-driver"`
	ModelStrict             bool     `mapstructure:"model-strict"`
	ModelIgnoreColumns      []string `mapstructure:"model-ignore-columns"`
	ModelIgnoreColumnsTable []struct {
		Table   string   `mapstructure:"table"`
		Columns []string `mapstructure:"columns"`
	} `mapstructure:"model-ignore-columns-table"`
	ModelSchema           string   `mapstructure:"model-schema"`
	ModelDatasource       bool     `mapstructure:"model-datasource"`
	ModelDatasourceUrl    []string `mapstructure:"model-datasource-url"`
	ModelDatasourceTable  []string `mapstructure:"model-datasource-table"`
	ModelCache            bool     `mapstructure:"model-cache"`
	ModelCacheTable       []string `mapstructure:"model-cache-table"`
	ModelCachePrefix      string   `mapstructure:"model-cache-prefix"`
	ModelNewOriginal      bool     `mapstructure:"model-new-original"`
	ModelCacheExpiryTable []struct {
		Table          string `mapstructure:"table"`
		Expiry         int64  `mapstructure:"expiry"`
		NotFoundExpiry int64  `mapstructure:"not-found-expiry"`
	} `mapstructure:"model-cache-expiry-table"`
	MongoType        []string `mapstructure:"mongo-type"`
	MongoCache       bool     `mapstructure:"mongo-cache"`
	MongoCachePrefix string   `mapstructure:"mongo-cache-prefix"`
	MongoCacheType   []string `mapstructure:"mongo-cache-type"`

	// Gen Sub Command
	Swagger GenSwaggerConfig `mapstructure:"swagger"`

	Zrpcclient GenZrpcclientConfig `mapstructure:"zrpcclient"`
}

type SkillsConfig struct {
	Init SkillsInitConfig `mapstructure:"init"`
}

type SkillsInitConfig struct {
	Output string `mapstructure:"output"`
}

type GenSwaggerConfig struct {
	Desc       []string `mapstructure:"desc"`
	DescIgnore []string `mapstructure:"desc-ignore"`
	Output     string   `mapstructure:"output"`
	Route2Code bool     `mapstructure:"route2code"`
	Merge      bool     `mapstructure:"merge"`
}

type GenZrpcclientConfig struct {
	Hooks HooksConfig `mapstructure:"hooks"`

	Desc         []string `mapstructure:"desc"`
	DescIgnore   []string `mapstructure:"desc-ignore"`
	ProtoInclude []string `mapstructure:"proto-include"`
	Output       string   `mapstructure:"output"`
	GoVersion    string   `mapstructure:"goVersion"`
	GoModule     string   `mapstructure:"goModule"`
	GoPackage    string   `mapstructure:"goPackage"`
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
	Delete ServerlessDeleteConfig `mapstructure:"delete"`
}

type ServerlessDeleteConfig struct {
	Plugin []string `mapstructure:"plugin"`
}

type MigrateConfig struct {
	Driver             string `mapstructure:"driver"`
	DataSourceUrl      string `mapstructure:"datasource-url"`
	Source             string `mapstructure:"source"`
	SourceAppendDriver bool   `mapstructure:"source-append-driver"`
	XMigrationsTable   string `mapstructure:"x-migrations-table"`
}

type FormatConfig struct {
	GitChange   bool `mapstructure:"git-change"`
	DisplayDiff bool `mapstructure:"display-diff"`
}

type AddConfig struct {
	Output string `mapstructure:"output"`

	Api          AddApiConfig          `mapstructure:"api"`
	Proto        AddProtoConfig        `mapstructure:"proto"`
	Sql          AddSqlConfig          `mapstructure:"sql"`
	SqlMigration AddSqlMigrationConfig `mapstructure:"sql-migration"`
}

type AddApiConfig struct{}

type AddProtoConfig struct{}

type AddSqlConfig struct{}

type AddSqlMigrationConfig struct{}

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

var (
	ToolVersionOnce  sync.Once
	ToolVersionValue ToolVersion
)

type ToolVersion struct {
	ProtocVersion             *version.Version
	GoctlVersion              *version.Version
	ProtocGenGoVersion        *version.Version
	ProtocGenGoGrpcVersion    *version.Version
	ProtocGenOpenapiv2Version *version.Version
}

func (c *Config) ToolVersion() ToolVersion {
	ToolVersionOnce.Do(func() {
		goctlVersionResp, _ := execx.Run("goctl -v", "")
		logx.Debugf("goctl version: %s", goctlVersionResp)
		versionInfo := strings.Split(goctlVersionResp, " ")
		if len(versionInfo) >= 3 {
			ToolVersionValue.GoctlVersion, _ = version.NewVersion(strings.TrimPrefix(versionInfo[2], "v"))
		}

		protocVersionResp, _ := protoc.Version()
		logx.Debugf("protoc version: %s", protocVersionResp)
		ToolVersionValue.ProtocVersion, _ = version.NewVersion(strings.TrimPrefix(protocVersionResp, "libprotoc "))

		protocGenGoVersionResp, _ := protocgengo.Version()
		logx.Debugf("protoc-gen-go version: %s", protocGenGoVersionResp)
		ToolVersionValue.ProtocGenGoVersion, _ = version.NewVersion(strings.TrimPrefix(protocGenGoVersionResp, "v"))

		protocGenGoGrpcVersionResp, _ := protocgengogrpc.Version()
		logx.Debugf("protoc-gen-go-grpc version: %s", protocGenGoGrpcVersionResp)
		ToolVersionValue.ProtocGenGoGrpcVersion, _ = version.NewVersion(strings.TrimPrefix(protocGenGoGrpcVersionResp, "v"))

		protocGenOpenapiv2VersionResp, _ := execx.Run("protoc-gen-openapiv2 --version", "")
		logx.Debugf("protoc-gen-openapiv2 version: %s", protocGenOpenapiv2VersionResp)
		versionInfo = strings.Split(protocGenOpenapiv2VersionResp, " ")
		if len(versionInfo) >= 2 {
			ToolVersionValue.ProtocGenOpenapiv2Version, _ = version.NewVersion(strings.TrimSuffix(strings.TrimPrefix(versionInfo[1], "v"), ","))
		}
	})

	return ToolVersionValue
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
