package templateinit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/zeromicro/go-zero/core/color"
)

func Init(cc config.TemplateConfig) error {
	if cc.Init.Remote != "" && cc.Init.Branch != "" {
		_ = os.MkdirAll(cc.Init.Output, 0o755)
		fmt.Printf("%s templates into '%s/templates/%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), cc.Init.Output, cc.Init.Branch)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := git.PlainCloneContext(ctx, filepath.Join(cc.Init.Output), false, &git.CloneOptions{
			SingleBranch:  true,
			URL:           cc.Init.Remote,
			Depth:         0,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + cc.Init.Branch),
		})
		if err != nil {
			return err
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))

		return nil
	}

	err := embeded.WriteTemplateDir("", cc.Init.Output)
	return err
}
