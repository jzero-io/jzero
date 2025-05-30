package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"{{ .Module }}/internal/config"
)

func (sc *ServiceContext) GetConfig() (config.Config, error) {
	return sc.Config.GetConfig()
}

func (sc *ServiceContext) MustGetConfig() config.Config {
	c, err := sc.GetConfig()
	logx.Must(err)
	return c
}

func (sc *ServiceContext) SetConfigListener() {
	sc.Config.AddListener(func() {
	    v, err := sc.GetConfig()
		if err != nil {
		    logx.Errorf("reload config error: %v", err)
		    return
		}

		logx.Infof("reload config successfully")
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

        // add custom logic here
	})
}