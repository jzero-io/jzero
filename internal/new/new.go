package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
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
}

type JzeroNew struct {
	TemplateData map[string]interface{}
	AppDir       string
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

	// template, register global data
	templateData := map[string]interface{}{
		"Module": Module,
		"APP":    AppName,
		"AppDir": AppDir,
	}
	jn := JzeroNew{
		TemplateData: templateData,
		AppDir:       AppDir,
	}

	err = jn.New(filepath.Join("app"))
	if err != nil {
		return err
	}
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

func (jn *JzeroNew) New(dirname string) error {
	dir := embeded.ReadTemplateDir(dirname)
	for _, file := range dir {
		if file.IsDir() {
			err := jn.New(filepath.Join(dirname, file.Name()))
			if err != nil {
				return err
			}
		}
		fileBytes, err := templatex.ParseTemplate(jn.TemplateData, embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())))
		if err != nil {
			return err
		}
		filename := strings.TrimSuffix(file.Name(), ".tpl")
		rel, err := filepath.Rel(filepath.Join("app"), filepath.Join(dirname, filename))
		if err != nil {
			return err
		}
		path := filepath.Join(jn.AppDir, rel)
		err = checkWrite(filepath.Join(Output, path), fileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
