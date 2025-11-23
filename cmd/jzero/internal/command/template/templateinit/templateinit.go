package templateinit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/zeromicro/go-zero/core/color"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
)

func Run() error {
	if config.C.Template.Init.Remote != "" && config.C.Template.Init.Branch != "" {
		target := filepath.Join(config.C.Template.Init.Output, config.C.Template.Init.Branch)
		_ = os.MkdirAll(target, 0o755)
		fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), target)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := git.PlainCloneContext(ctx, target, false, &git.CloneOptions{
			SingleBranch:  true,
			URL:           config.C.Template.Init.Remote,
			Depth:         0,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + config.C.Template.Init.Branch),
		})
		if err != nil {
			return err
		}
		_ = os.RemoveAll(filepath.Join(target, ".git"))
		fmt.Println(color.WithColor("Done", color.FgGreen))
		return nil
	}
	fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Initializing embedded", color.FgGreen), config.C.Template.Init.Output)
	err := embeded.WriteTemplateDir("", config.C.Template.Init.Output)
	if err != nil {
		return err
	}
	fmt.Println(color.WithColor("Done", color.FgGreen))
	return nil
}
