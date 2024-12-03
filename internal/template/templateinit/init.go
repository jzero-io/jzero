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

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
)

func Run(c config.Config) error {
	if c.Template.Init.Remote != "" && c.Template.Init.Branch != "" {
		target := filepath.Join(c.Template.Init.Output, c.Template.Init.Branch)
		_ = os.MkdirAll(target, 0o755)
		fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), target)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := git.PlainCloneContext(ctx, target, false, &git.CloneOptions{
			SingleBranch:  true,
			URL:           c.Template.Init.Remote,
			Depth:         0,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + c.Template.Init.Branch),
		})
		if err != nil {
			return err
		}
		_ = os.RemoveAll(filepath.Join(target, ".git"))
		fmt.Println(color.WithColor("Done", color.FgGreen))
		return nil
	}
	fmt.Printf("%s templates into '%s', please wait...\n", color.WithColor("Initializing embedded", color.FgGreen), c.Template.Init.Output)
	err := embeded.WriteTemplateDir("", c.Template.Init.Output)
	if err != nil {
		return err
	}
	fmt.Println(color.WithColor("Done", color.FgGreen))
	return nil
}
