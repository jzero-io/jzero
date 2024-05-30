package cmd

import (
    "fmt"

	"github.com/jzero-io/jzero-contrib/swaggerv2"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/handler"
	"{{ .Module }}/internal/svc"
	"{{ .Module }}/internal/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "{{ .APP }} server",
	Long:  "{{ .APP }} server",
	Run: func(cmd *cobra.Command, args []string) {
		Start(cfgFile)
	},
}

func Start(cfgFile string) {
	var c config.Config
	conf.MustLoad(cfgFile, &c)
    // set up logger
    if err := logx.SetUp(c.Log.LogConf); err != nil {
        logx.Must(err)
    }

	ctx := svc.NewServiceContext(c)
	start(ctx)
}

func start(ctx *svc.ServiceContext) {
	s := server.RegisterZrpc(ctx.Config, ctx)
	gw := gateway.MustNewServer(ctx.Config.Gateway.GatewayConf)

    // gw add api routes
    handler.RegisterHandlers(gw.Server, ctx)

    // gw add swagger routes. If you do not want it, you can delete this line
    swaggerv2.RegisterRoutes(gw.Server)

	// gw add routes
    // You can use gw.Server.AddRoutes() to add your own handler
    // for example: add a func handler.RegisterMyHandlers() in this line on handler dir

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

    fmt.Printf("Starting rpc server at %s...\n", ctx.Config.Zrpc.ListenOn)
	fmt.Printf("Starting gateway server at %s:%d...\n", ctx.Config.Gateway.Host, ctx.Config.Gateway.Port)
	group.Start()

}

func init() {
	rootCmd.AddCommand(serverCmd)
}
