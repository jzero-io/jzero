package serverlessnew

import (
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/new"
	"github.com/jzero-io/jzero/internal/serverless/serverlessbuild"
)

func Run(args []string) error {
	config.C.New = config.NewConfig{
		Home:     config.C.Serverless.New.Home,
		Module:   config.C.Serverless.New.Module,
		Output:   filepath.Join("plugins", args[0]),
		Remote:   config.C.Serverless.New.Remote,
		Frame:    config.C.Serverless.New.Frame,
		Branch:   config.C.Serverless.New.Branch,
		Local:    config.C.Serverless.New.Local,
		Style:    config.C.Serverless.New.Style,
		Features: []string{"serverless"},
	}
	if config.C.Serverless.New.Core {
		config.C.New.Features = append(config.C.New.Features, "serverless_core")
		config.C.New.Output = args[0]
	}
	err := new.Run(config.C, args[0])
	if err != nil {
		return err
	}

	if config.C.Serverless.New.Core {
		return nil
	}
	path := "./" + "plugins/" + args[0]

	if _, err := os.Stat("go.work"); err == nil {
		goWork, _ := os.ReadFile("go.work")
		work, err := modfile.ParseWork("", goWork, nil)
		if err != nil {
			return err
		}
		err = work.AddUse(path, "")
		if err != nil {
			return err
		}
		if err = os.WriteFile("go.work", modfile.Format(work.Syntax), 0o644); err != nil {
			return err
		}
	}
	return serverlessbuild.Run()
}
