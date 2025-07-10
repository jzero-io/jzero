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
	"{{ .Module }}/internal/global"
	"{{ .Module }}/internal/middleware"
	"{{ .Module }}/internal/svc"
	"{{ .Module }}/internal/server"
	"{{ .Module }}/internal/custom"
	{{ if not .Serverless }}"{{ .Module }}/plugins"{{end}}
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

    	// print banner
        printBanner(c)
        // print version
        printVersion()

    	svcCtx := svc.NewServiceContext(cc)
    	global.ServiceContext = *svcCtx
    	run(svcCtx)
	},
}

func run(svcCtx *svc.ServiceContext) {
    var err error
    svcCtx.MustGetConfig().Gateway.Upstreams[0].ProtoSets, err = pb.WriteToLocal(pb.Embed)
    logx.Must(err)

	zrpcServer := zrpc.MustNewServer(svcCtx.MustGetConfig().Zrpc.RpcServerConf, func(grpcServer *grpc.Server) {
		server.RegisterZrpcServer(grpcServer, svcCtx)
        {{if not .Serverless }}// register plugins
        plugins.LoadPlugins(grpcServer, svcCtx){{end}}
		if svcCtx.MustGetConfig().Zrpc.Mode == service.DevMode || svcCtx.MustGetConfig().Zrpc.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	gatewayServer := gateway.MustNewServer(svcCtx.MustGetConfig().Gateway.GatewayConf, middleware.WithHeaderProcessor())

	ctm := custom.New(zrpcServer, gatewayServer)
	ctm.Init()

	// register middleware
	middleware.Register(zrpcServer, gatewayServer)

	group := service.NewServiceGroup()
	group.Add(zrpcServer)
	group.Add(gatewayServer)
	group.Add(ctm)

	logx.Infof("Starting rpc server at %s...", svcCtx.MustGetConfig().Zrpc.ListenOn)
	logx.Infof("Starting gateway server at %s:%d...", svcCtx.MustGetConfig().Gateway.Host, svcCtx.MustGetConfig().Gateway.Port)

	group.Start()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
