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
		_ = os.MkdirAll(c.Template.Init.Output, 0o755)
		fmt.Printf("%s templates into '%s/templates/%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), c.Template.Init.Output, c.Template.Init.Branch)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := git.PlainCloneContext(ctx, filepath.Join(c.Template.Init.Output), false, &git.CloneOptions{
			SingleBranch:  true,
			URL:           c.Template.Init.Remote,
			Depth:         0,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + c.Template.Init.Branch),
		})
		if err != nil {
			return err
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))

		return nil
	}

	err := embeded.WriteTemplateDir("", c.Template.Init.Output)
	return err
}
