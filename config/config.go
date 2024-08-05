package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type Config struct {
	Version string `mapstructure:"version"`

	Gen GenConfig `mapstructure:"gen"`
}

type GenConfig struct {
	Home                      string   `mapstructure:"home"`
	Style                     string   `mapstructure:"style"`
	ChangeReplaceTypes        bool     `mapstructure:"change-replace-types"`
	RemoveSuffix              bool     `mapstructure:"remove-suffix"`
	ModelMysqlIgnoreColumns   []string `mapstructure:"model-mysql-ignore-columns"`
	ModelMysqlDatasource      bool     `mapstructure:"model-mysql-datasource"`
	ModelMysqlDatasourceUrl   string   `mapstructure:"model-mysql-datasource-url"`
	ModelMysqlDatasourceTable []string `mapstructure:"model-mysql-datasource-table"`
	ModelMysqlCache           bool     `mapstructure:"model-mysql-cache"`
	ModelMysqlCachePrefix     string   `mapstructure:"model-mysql-cache-prefix"`
}

func SetConfig(cfgFile string, command string, flagSet *pflag.FlagSet) (*Config, error) {
	v := viper.New()
	if err := v.BindPFlags(flagSet); err != nil {
		return nil, err
	}
	viper.Set(command, v.AllSettings())

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	if pathx.FileExists(cfgFile) {
		x := viper.New()
		x.SetConfigFile(cfgFile)
		if err := x.ReadInConfig(); err != nil {
			return nil, err
		}
		if err := x.Unmarshal(&c); err != nil {
			return nil, err
		}
	}
	logx.Debugf("config file version: %s", c.Version)
	return &c, nil
}
