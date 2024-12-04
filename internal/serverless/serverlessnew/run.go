package serverlessnew

import (
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/new"
)

func Run(args []string) error {
	if config.C.Serverless.New.Core {
		config.C.Serverless.New.Features = append(config.C.Serverless.New.Features, "serverless_core")
	} else {
		config.C.Serverless.New.Features = append(config.C.Serverless.New.Features, "serverless")
	}
	config.C.New = config.NewConfig{
		Home:     config.C.Serverless.New.Home,
		Module:   config.C.Serverless.New.Module,
		Output:   filepath.Join("plugins", args[0]),
		Remote:   config.C.Serverless.New.Remote,
		Frame:    config.C.Serverless.New.Frame,
		Branch:   config.C.Serverless.New.Branch,
		Local:    config.C.Serverless.New.Local,
		Style:    config.C.Serverless.New.Style,
		Features: config.C.Serverless.New.Features,
	}
	if config.C.Serverless.New.Core {
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
	return nil
}
