package cmd

import (
    "os"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/common-nighthawk/go-figure"
	"golang.org/x/sync/errgroup"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/handler"
	"{{ .Module }}/internal/svc"
	"{{ .Module }}/internal/middleware"
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
	config.C = c

    // set up logger
    if err := logx.SetUp(c.Log.LogConf); err != nil {
        logx.Must(err)
    }
    if c.Log.LogConf.Mode != "console" {
        logx.AddWriter(logx.NewWriter(os.Stdout))
    }

	ctx := svc.NewServiceContext(c)
	start(ctx)
}

func start(svcCtx *svc.ServiceContext) {
	server := rest.MustNewServer(svcCtx.Config.Rest.RestConf)
	middleware.Register(server)

	// server add api handlers
	handler.RegisterHandlers(server, svcCtx)

	// server add custom routes
    svcCtx.Custom.AddRoutes(server)

	group := service.NewServiceGroup()
	group.Add(server)

	// shutdown listener
    waitExit := proc.AddShutdownListener(svcCtx.Custom.Stop)

    eg := errgroup.Group{}
    eg.Go(func() error {
    	printBanner(svcCtx.Config)
    	logx.Infof("Starting rest server at %s:%d...", svcCtx.Config.Rest.Host, svcCtx.Config.Rest.Port)
    	group.Start()
    	return nil
    })

    // add custom start logic
    eg.Go(func() error {
    	svcCtx.Custom.Start()
    	return nil
    })

    if err := eg.Wait(); err != nil {
    	panic(err)
    }

    waitExit()
}

func printBanner(c config.Config) {
	figure.NewColorFigure(c.Banner.Text, c.Banner.FontName, c.Banner.Color, true).Print()
}


func init() {
	rootCmd.AddCommand(serverCmd)
}
