package genswagger

import (
	"fmt"
	"github.com/jzero-io/jzero/cmd/gen"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"os"
	"path/filepath"
)

func Gen(_ *cobra.Command, _ []string) error {
	wd, _ := os.Getwd()
	apiDirName := filepath.Join(wd, "app", "desc", "api")
	mainApiFile := gen.GetMainApiFilePath(apiDirName)
	defer os.Remove(mainApiFile)
	if mainApiFile != "" {
		command := fmt.Sprintf("goctl api plugin -plugin goctl-swagger=\"swagger -filename swagger.json --schemes http --host 127.0.0.1:8001\" -api %s -dir .", mainApiFile)
		_, err := execx.Run(command, wd)
		if err != nil {
			return err
		}
	}
	return nil
}
