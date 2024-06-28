package templateinit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
)

var (
	Home   string
	Remote string
	Branch string
)

func Init(_ *cobra.Command, _ []string) error {
	if Remote != "" && Branch != "" {
		_ = os.MkdirAll(Home, 0o755)
		fmt.Printf("%s templates into '%s/templates/%s', please wait...\n", color.WithColor("Cloning", color.FgGreen), Home, Branch)
		_, err := git.PlainClone(filepath.Join(Home), false, &git.CloneOptions{
			SingleBranch:  true,
			URL:           Remote,
			Depth:         0,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + Branch),
		})
		if err != nil {
			return err
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))

		return nil
	}

	err := embeded.WriteTemplateDir("", Home)
	return err
}
