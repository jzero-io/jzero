package new

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/embeded"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Module  string
	Output  string
	AppDir  string
	AppName string
	Remote  string
	Branch  string

	Version string
)

type TemplateData struct {
	Module string
	APP    string
	AppDir string
}

func NewProject(_ *cobra.Command, _ []string) error {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if embeded.Home == "" {
		embeded.Home = filepath.Join(homeDir, ".jzero", Version)
	}

	err = os.MkdirAll(filepath.Join(Output, AppDir), 0o755)
	cobra.CheckErr(err)

	_, err = execx.Run(fmt.Sprintf("go mod init %s", Module), filepath.Join(Output, AppDir))
	cobra.CheckErr(err)

	templateData := TemplateData{
		Module: Module,
		APP:    AppName,
		AppDir: AppDir,
	}

	jzeroRoot := JzeroRoot{TemplateData: templateData, AppDir: AppDir}
	err = jzeroRoot.New()
	cobra.CheckErr(err)

	jzeroEtc := JzeroEtc{TemplateData: templateData, AppDir: AppDir}
	err = jzeroEtc.New()
	cobra.CheckErr(err)

	jzeroCmd := JzeroCmd{TemplateData: templateData, AppDir: AppDir}
	err = jzeroCmd.New()
	cobra.CheckErr(err)

	jzeroProto := JzeroProto{TemplateData: templateData, AppDir: AppDir}
	err = jzeroProto.New()
	cobra.CheckErr(err)

	jzeroApi := JzeroApi{TemplateData: templateData, AppDir: AppDir}
	err = jzeroApi.New()
	cobra.CheckErr(err)

	return nil
}

func checkWrite(path string, bytes []byte) error {
	var err error
	if len(bytes) == 0 {
		return nil
	}
	if !pathx.FileExists(filepath.Dir(path)) {
		err = os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			return err
		}
	}

	bytesFormat := bytes
	if filepath.Ext(path) == ".go" {
		bytesFormat, err = gosimports.Process("", bytes, &gosimports.Options{FormatOnly: true})
		if err != nil {
			return err
		}
	}

	return os.WriteFile(path, bytesFormat, 0o644)
}
