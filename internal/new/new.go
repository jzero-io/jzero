package new

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type TemplateData struct {
	Module string
	APP    string
}

type JzeroNew struct {
	TemplateData map[string]interface{}
	nc           config.NewConfig
}

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

func NewProject(nc config.NewConfig, appName string) error {
	err := os.MkdirAll(nc.Output, 0o755)
	cobra.CheckErr(err)

	templateData, err := NewTemplateData()
	cobra.CheckErr(err)
	templateData["Features"] = nc.Features
	templateData["Module"] = nc.Module
	templateData["APP"] = appName

	jn := JzeroNew{
		TemplateData: templateData,
		nc:           nc,
	}

	gfs, err := jn.New(filepath.Join("app"))
	if err != nil {
		return err
	}

	for _, gf := range gfs {
		if !gf.Skip {
			err = checkWrite(gf.Path, gf.Content.Bytes())
			if err != nil {
				return err
			}
		}
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
			return errors.Wrapf(err, "format %s", path)
		}
	}

	return os.WriteFile(path, bytesFormat, 0o644)
}

func (jn *JzeroNew) New(dirname string) ([]*GeneratedFile, error) {
	var gsf []*GeneratedFile

	dir := embeded.ReadTemplateDir(dirname)
	for _, file := range dir {
		if file.IsDir() {
			files, err := jn.New(filepath.Join(dirname, file.Name()))
			if err != nil {
				return nil, err
			}
			gsf = append(gsf, files...)
		}
		filename := strings.TrimSuffix(file.Name(), ".tpl")
		rel, err := filepath.Rel(filepath.Join("app"), filepath.Join(dirname, filename))
		if err != nil {
			return nil, err
		}
		fileBytes, err := templatex.ParseTemplate(jn.TemplateData, embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())))
		if err != nil {
			return nil, errors.Wrapf(err, "parse template: %s", filepath.Join(dirname, file.Name()))
		}

		stylePath := filepath.Join(filepath.Dir(rel), filename)
		if filepath.Ext(filename) == ".go" {
			formatFilename, err := format.FileNamingFormat(jn.nc.Style, filename[0:len(filename)-len(filepath.Ext(filename))])
			if err != nil {
				return nil, err
			}
			stylePath = filepath.Join(filepath.Dir(rel), formatFilename+filepath.Ext(filename))
		}

		// specify
		if filename == "go.mod" && jn.nc.SubModule {
			continue
		}

		gsf = append(gsf, &GeneratedFile{
			Path:    filepath.Join(jn.nc.Output, stylePath),
			Content: *bytes.NewBuffer(fileBytes),
			// Because this is a special directory for jzero
			// It is deleted to support generating all server code under the premise of only desc directory
			// Or if this file has been existed, just ignore write
			Skip: pathx.FileExists(filepath.Join(jn.nc.Output, "desc")) && strings.HasPrefix(rel, "desc"),
		})
	}
	return gsf, nil
}
