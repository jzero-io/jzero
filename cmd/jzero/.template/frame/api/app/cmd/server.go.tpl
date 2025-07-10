package cmd

import (
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/spf13/cobra"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/custom"
	"{{ .Module }}/internal/global"
	"{{ .Module }}/internal/middleware"
	"{{ .Module }}/internal/handler"
	"{{ .Module }}/internal/svc"
	{{ if not .Serverless }}"{{ .Module }}/plugins"{{end}}
)

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
		if err = logx.SetUp(c.Log.LogConf); err != nil {
			logx.Must(err)
		}

		// print banner
		printBanner(c)
		// print version
		printVersion()

		svcCtx := svc.NewServiceContext(cc)
		svcCtx.Middleware = middleware.NewMiddleware()
		global.ServiceContext = *svcCtx
		run(svcCtx)
	},
}

func run(svcCtx *svc.ServiceContext) {
    server := rest.MustNewServer(svcCtx.MustGetConfig().Rest.RestConf)

	ctm := custom.New(server)
	ctm.Init()

	handler.RegisterHandlers(server, svcCtx)
	middleware.Register(server)

	{{ if not .Serverless }}// load plugins
	plugins.LoadPlugins(server, svcCtx){{end}}

	group := service.NewServiceGroup()
	group.Add(server)
	group.Add(ctm)

	group.Start()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
