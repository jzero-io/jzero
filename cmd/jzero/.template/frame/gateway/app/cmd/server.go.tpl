package cmd

import (
    "github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/swaggerv2"
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
		cc := configcenter.MustNewConfigCenter[config.Config](configcenter.Config{
			Type: "yaml",
		}, subscriber.MustNewFsnotifySubscriber(cmd.Flag("config").Value.String(), subscriber.WithUseEnv(true)))

        // set up logger
    	logx.Must(logx.SetUp(cc.MustGetConfig().Log.LogConf))

    	// print banner
        printBanner(cc.MustGetConfig().Banner)
        // print version
        printVersion()

        // create service context
    	svcCtx := svc.NewServiceContext(cc)

    	var err error
    	// write protosets to local
        cc.MustGetConfig().Gateway.Upstreams[0].ProtoSets, err = pb.WriteToLocal(pb.Embed)
        logx.Must(err)

        // create zrpc server
        zrpcServer := zrpc.MustNewServer(cc.MustGetConfig().Zrpc.RpcServerConf, func(grpcServer *grpc.Server) {
        	server.RegisterZrpcServer(grpcServer, svcCtx)
               {{if not .Serverless }}// register plugins
               plugins.LoadPlugins(grpcServer, svcCtx){{end}}
        	if cc.MustGetConfig().Zrpc.Mode == service.DevMode || cc.MustGetConfig().Zrpc.Mode == service.TestMode {
        		reflection.Register(grpcServer)
        	}
        })
        // create gateway server
        gatewayServer := gateway.MustNewServer(cc.MustGetConfig().Gateway.GatewayConf, middleware.WithHeaderProcessor())
        // register swagger routes
        swaggerv2.RegisterRoutes(gatewayServer.Server)
        // // create custom server
        customServer := custom.New()

        // register middleware
        middleware.Register(zrpcServer, gatewayServer)

        group := service.NewServiceGroup()
        group.Add(zrpcServer)
        group.Add(gatewayServer)
        group.Add(customServer)

        logx.Infof("Starting rpc server at %s...", cc.MustGetConfig().Zrpc.ListenOn)
        logx.Infof("Starting gateway server at %s:%d...", cc.MustGetConfig().Gateway.Host, cc.MustGetConfig().Gateway.Port)
        group.Start()
	},
}

func printBanner(c config.BannerConf) {
	figure.NewColorFigure(c.Text, c.FontName, c.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
