package new

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"
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
	TemplateData map[string]any
	nc           config.NewConfig
	base         string
}

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

func Run(c config.Config, appName string) error {
	err := os.MkdirAll(c.New.Output, 0o755)
	cobra.CheckErr(err)

	templateData, err := NewTemplateData()
	cobra.CheckErr(err)
	templateData["Features"] = c.New.Features
	templateData["Module"] = c.New.Module
	templateData["APP"] = appName
	if abs, err := filepath.Abs(c.New.Output); err == nil {
		templateData["DirName"] = filepath.Base(abs)
	}
	var base string
	if c.New.Frame != "" {
		// 使用内置 frame
		base = filepath.Join("frame", c.New.Frame, "app")
	} else {
		// 使用第三方模板仓库的模板
		base = filepath.Join("app")
	}

	jn := JzeroNew{
		TemplateData: templateData,
		nc:           c.New,
		base:         base,
	}

	gfs, err := jn.New(base)
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
		rel, err := filepath.Rel(jn.base, filepath.Join(dirname, filename))
		if err != nil {
			return nil, err
		}
		fileBytes, err := templatex.ParseTemplate(jn.TemplateData, embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())))
		if err != nil {
			fileBytes = embeded.ReadTemplateFile(filepath.Join(dirname, file.Name()))
		}

		stylePath := filepath.Join(filepath.Dir(rel), filename)

		// specify
		if filename == "go.mod" && jn.nc.Mono {
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
