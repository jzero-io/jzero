package genswagger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jzero-io/jzero/internal/gen"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Dir      string
	ApiDir   string
	ProtoDir string
)

func Gen(_ *cobra.Command, _ []string) error {
	wd, _ := os.Getwd()

	mainApiFile := gen.GetMainApiFilePath(ApiDir)
	if mainApiFile != filepath.Join("desc", "api", "main.api") {
		defer os.Remove(mainApiFile)
	}

	if !pathx.FileExists(Dir) {
		_ = os.MkdirAll(Dir, 0o755)
	}

	// gen swagger by desc/api
	if mainApiFile != "" {
		apiFile := fmt.Sprintf("%s.swagger.json", gen.GetApiServiceName(ApiDir))
		cmd := exec.Command("goctl", "api", "plugin", "-plugin", "goctl-swagger=swagger -filename "+apiFile+" --schemes http", "-api", mainApiFile, "-dir", Dir)
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Wrap(err, string(resp))
		}
		fmt.Println(string(resp))
	}

	if pathx.FileExists(ProtoDir) {
		protoDirFile, err := os.ReadDir(ProtoDir)
		if err != nil {
			return err
		}
		for _, protoFile := range protoDirFile {
			if protoFile.IsDir() {
				continue
			}
			if filepath.Ext(protoFile.Name()) == ".proto" {
				command := fmt.Sprintf("protoc -I%s %s --openapiv2_out=%s", ProtoDir, filepath.Join(ProtoDir, protoFile.Name()), Dir)
				_, err := execx.Run(command, wd)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
