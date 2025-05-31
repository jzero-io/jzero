package cmd

import (
    "os"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/spf13/cobra"
    "github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/common-nighthawk/go-figure"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/middleware"
	"{{ .Module }}/internal/server"
	"{{ .Module }}/internal/svc"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "{{ .APP }} server",
	Long:  "{{ .APP }} server",
	Run: func(cmd *cobra.Command, args []string) {
		cc := configurator.MustNewConfigCenter[config.Config](configurator.Config{
			Type: "yaml",
		}, subscriber.MustNewFsnotifySubscriber(cfgFile, subscriber.WithUseEnv(true)))
		c, err := cc.GetConfig()
		logx.Must(err)

        // set up logger
        if err := logx.SetUp(c.Log.LogConf); err != nil {
            logx.Must(err)
        }
    	if c.Log.LogConf.Mode != "console" {
    	    logx.AddWriter(logx.NewWriter(os.Stdout))
    	}

    	svcCtx := svc.NewServiceContext(cc)
    	run(svcCtx)
	},
}

func run(svcCtx *svc.ServiceContext) {
    c := svcCtx.MustGetConfig()

	zrpc := server.RegisterZrpc(c, svcCtx)
    middleware.Register(zrpc)

	group := service.NewServiceGroup()
	group.Add(zrpc)
	group.Add(svcCtx.Custom)

	printBanner(c)
    logx.Infof("Starting rpc server at %s...", c.Zrpc.ListenOn)
    group.Start()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}


func init() {
	rootCmd.AddCommand(serverCmd)
}
