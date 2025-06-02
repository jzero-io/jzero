package cmd

import (
    "path/filepath"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
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
		cc := configurator.MustNewConfigCenter[config.Config](configurator.Config{
			Type: "yaml",
		}, subscriber.MustNewFsnotifySubscriber(cfgFile, subscriber.WithUseEnv(true)))
		c, err := cc.GetConfig()
		logx.Must(err)

        // set up logger
    	logx.Must(logx.SetUp(c.Log.LogConf))

    	// write pb to local
        c.Gateway.Upstreams[0].ProtoSets, err = pb.WriteToLocal(pb.Embed, pb.WithFileMatchFunc(func(path string) bool {
			return filepath.Ext(path) == ".pb"
		}))
        logx.Must(err)

    	svcCtx := svc.NewServiceContext(cc)
    	run(svcCtx)
	},
}

func run(svcCtx *svc.ServiceContext) {
    c := svcCtx.MustGetConfig()

	zrpc := server.RegisterZrpc(c, svcCtx)
	gw := gateway.MustNewServer(c.Gateway.GatewayConf, middleware.WithHeaderProcessor())

	// register middleware
	middleware.Register(zrpc, gw)

	// gw add custom routes
	svcCtx.Custom.AddRoutes(gw)

	group := service.NewServiceGroup()
	group.Add(zrpc)
	group.Add(gw)
	group.Add(svcCtx.Custom)

	printBanner(c)
	printVersion()

	logx.Infof("Starting rpc server at %s...", c.Zrpc.ListenOn)
	logx.Infof("Starting gateway server at %s:%d...", c.Gateway.Host, c.Gateway.Port)

	group.Start()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
