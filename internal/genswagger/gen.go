package genswagger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/internal/gen"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var Dir string

func Gen(_ *cobra.Command, _ []string) error {
	wd, _ := os.Getwd()

	apiDirName := filepath.Join(wd, "app", "desc", "api")
	protoDirName := filepath.Join(wd, "app", "desc", "proto")

	mainApiFile := gen.GetMainApiFilePath(apiDirName)
	defer os.Remove(mainApiFile)
	if !pathx.FileExists(Dir) {
		_ = os.MkdirAll(Dir, 0o755)
	}

	// gen swagger by app/desc/api
	if mainApiFile != "" {
		command := fmt.Sprintf("goctl api plugin -plugin goctl-swagger=\"swagger -filename %s.swagger.json --schemes http\" -api %s -dir %s", gen.GetApiServiceName(apiDirName), mainApiFile, Dir)
		_, err := execx.Run(command, wd)
		if err != nil {
			return err
		}
	}

	// gen swagger by app/desc/proto
	if pathx.FileExists(protoDirName) {
		protoDir, err := os.ReadDir(protoDirName)
		if err != nil {
			return err
		}
		for _, protoFile := range protoDir {
			if protoFile.IsDir() {
				continue
			}
			if filepath.Ext(protoFile.Name()) == ".proto" {
				command := fmt.Sprintf("protoc -I./app/desc/proto ./app/desc/proto/%s --openapiv2_out=./app/desc/swagger", protoFile.Name())
				_, err := execx.Run(command, wd)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
