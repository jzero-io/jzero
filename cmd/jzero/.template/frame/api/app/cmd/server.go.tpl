package cmd

import (
	"github.com/spf13/cobra"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
    "github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/common-nighthawk/go-figure"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/middleware"
	"{{ .Module }}/internal/handler"
	"{{ .Module }}/internal/svc"
	{{ if has "serverless_core" .Features }}"{{ .Module }}/plugins"{{end}}
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

        svcCtx := svc.NewServiceContext(cc)
        run(svcCtx)
	},
}

func run(svcCtx *svc.ServiceContext) {
    c := svcCtx.MustGetConfig()

	server := rest.MustNewServer(c.Rest.RestConf)
	middleware.Register(server)

	// server add api handlers
	handler.RegisterHandlers(server, svcCtx)

	// server add custom routes
    svcCtx.Custom.AddRoutes(server)

    {{ if has "serverless_core" .Features }}// load plugins features
    plugins.LoadPlugins(server, *svcCtx){{end}}

	group := service.NewServiceGroup()
	group.Add(server)
	group.Add(svcCtx.Custom)

	printBanner(c)
	printVersion()

    logx.Infof("Starting rest server at %s:%d...", c.Rest.Host, c.Rest.Port)
    group.Start()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
