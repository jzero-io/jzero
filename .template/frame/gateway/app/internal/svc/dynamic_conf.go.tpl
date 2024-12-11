package svc

import (
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"

	"{{ .Module }}/internal/config"
)

func (sc *ServiceContext) DynamicConfListener(cc configurator.Configurator[config.Config]) {
	cc.AddListener(func() {
		logLevel := sc.Config.Log.Level
		logx.Infof("config file changed")
		if v, err := cc.GetConfig(); err == nil {
			if v.Log.Level != logLevel {
				logx.Infof("log level changed: %s", v.Log.Level)
				switch v.Log.Level {
				case "debug":
					logx.SetLevel(logx.DebugLevel)
				case "info":
					logx.SetLevel(logx.InfoLevel)
				case "error":
					logx.SetLevel(logx.ErrorLevel)
				case "severe":
					logx.SetLevel(logx.SevereLevel)
				}
			}

			config.C = v
			sc.Config = v
		}
	})
}
