package cmd

import (
    "os"

	"github.com/common-nighthawk/go-figure"
	"github.com/jzero-io/jzero-contrib/embedx"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"

	"{{ .Module }}/desc/pb"
	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/middleware"
	"{{ .Module }}/internal/svc"
	"{{ .Module }}/internal/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "{{ .APP }} server",
	Long:  "{{ .APP }} server",
	Run: func(cmd *cobra.Command, args []string) {
	    var c config.Config
    	conf.MustLoad(cfgFile, &c)
    	config.C = c

    	// write pb to local
        var err error
        c.Gateway.Upstreams[0].ProtoSets, err = embedx.WriteToLocalTemp(pb.Embed, embedx.WithFileMatchFunc(func(path string) bool {
			return filepath.Ext(path) == ".pb"
		}))
        if err != nil {
        	logx.Must(err)
        }

    	// set up logger
    	if err = logx.SetUp(c.Log.LogConf); err != nil {
    		logx.Must(err)
    	}
    	if c.Log.LogConf.Mode != "console" {
            logx.AddWriter(logx.NewWriter(os.Stdout))
        }

    	ctx := svc.NewServiceContext(c)
    	run(ctx)
	},
}

func run(svcCtx *svc.ServiceContext) {
	zrpc := server.RegisterZrpc(svcCtx.Config, svcCtx)
	gw := gateway.MustNewServer(svcCtx.Config.Gateway.GatewayConf, middleware.WithHeaderProcessor())

	// register middleware
	middleware.Register(zrpc, gw)

	// gw add custom routes
	svcCtx.Custom.AddRoutes(gw)

	group := service.NewServiceGroup()
	group.Add(zrpc)
	group.Add(gw)
	group.Add(svcCtx.Custom)

	printBanner(svcCtx.Config)
	logx.Infof("Starting rpc server at %s...", svcCtx.Config.Zrpc.ListenOn)
	logx.Infof("Starting gateway server at %s:%d...", svcCtx.Config.Gateway.Host, svcCtx.Config.Gateway.Port)

	group.Start()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
