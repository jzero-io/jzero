package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"{{ .Module }}/internal/config"
)

func (sc *ServiceContext) GetConfig() (config.Config, error) {
	return sc.ConfigCenter.GetConfig()
}

func (sc *ServiceContext) MustGetConfig() config.Config {
	c, err := sc.ConfigCenter.GetConfig()
	logx.Must(err)
	return c
}

func (svcCtx *ServiceContext) SetConfigListener() {
	svcCtx.ConfigCenter.AddListener(func() {
	    c, err := svcCtx.ConfigCenter.GetConfig()
		if err != nil {
		    logx.Errorf("reload config error: %v", err)
		    return
		}

		logx.Infof("reload config successfully")
		switch c.Log.Level {
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
