package new

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Module  string
	Output  string
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
}

func NewProject(_ *cobra.Command, _ []string) error {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if embeded.Home == "" {
		embeded.Home = filepath.Join(homeDir, ".jzero", Version)
	}

	err = os.MkdirAll(Output, 0o755)
	cobra.CheckErr(err)

	templateData, err := newTemplateData()
	cobra.CheckErr(err)

	jn := JzeroNew{
		TemplateData: templateData,
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
		bytesFormat, err = gosimports.Process("", bytes, &gosimports.Options{FormatOnly: true, Comments: true})
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
		filename := strings.TrimSuffix(file.Name(), ".tpl")
		rel, err := filepath.Rel(filepath.Join("app"), filepath.Join(dirname, filename))
		if err != nil {
			return err
		}
		fileBytes, err := templatex.ParseTemplate(jn.TemplateData, embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())))
		if err != nil {
			console.Warning("parse template file [%s] error, simple replace it", rel)
			fileBytes = bytes.ReplaceAll(embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())), []byte("{{ .Module }}"), []byte(cast.ToString(jn.TemplateData["Module"])))
		}

		path := filepath.Join(rel)
		err = checkWrite(filepath.Join(Output, path), fileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
