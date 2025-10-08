package cmd

import (
    "github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/spf13/cobra"
    "github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/common-nighthawk/go-figure"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/custom"
	"{{ .Module }}/internal/middleware"
	"{{ .Module }}/internal/server"
	"{{ .Module }}/internal/svc"
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

	    printBanner(cc.MustGetConfig().Banner)
	    printVersion()

        // create service context
    	svcCtx := svc.NewServiceContext(cc)
        // create zrpc server
	    zrpcServer := zrpc.MustNewServer(cc.MustGetConfig().Zrpc.RpcServerConf, func(grpcServer *grpc.Server) {
            server.RegisterZrpcServer(grpcServer, svcCtx)
                {{if not .Serverless }}// register plugins
                plugins.LoadPlugins(grpcServer, svcCtx){{end}}
            if cc.MustGetConfig().Zrpc.Mode == service.DevMode || cc.MustGetConfig().Zrpc.Mode == service.TestMode {
            	reflection.Register(grpcServer)
            }
        })
        // create custom server
	    customServer := custom.New()
        // register middleware
        middleware.Register(zrpcServer)

	    group := service.NewServiceGroup()
	    group.Add(zrpcServer)
	    group.Add(customServer)

        logx.Infof("Starting rpc server at %s...", cc.MustGetConfig().Zrpc.ListenOn)
        group.Start()
	},
}

func printBanner(c config.BannerConf) {
	figure.NewColorFigure(c.Text, c.FontName, c.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
