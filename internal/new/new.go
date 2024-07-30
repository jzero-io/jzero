package new

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Module       string
	Output       string
	AppName      string
	Remote       string
	Cache        bool
	Branch       string
	WithTemplate bool
	Style        string

	Features []string
)

type TemplateData struct {
	Module string
	APP    string
}

type JzeroNew struct {
	TemplateData map[string]interface{}
	Style        string
}

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

func NewProject(_ *cobra.Command, _ []string) error {
	err := os.MkdirAll(Output, 0o755)
	cobra.CheckErr(err)

	templateData, err := NewTemplateData(Features)
	cobra.CheckErr(err)

	jn := JzeroNew{
		TemplateData: templateData,
		Style:        Style,
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
			return err
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
			return nil, err
		}

		stylePath := filepath.Join(filepath.Dir(rel), filename)
		if filepath.Ext(filename) == ".go" {
			formatFilename, err := format.FileNamingFormat(jn.Style, filename[0:len(filename)-len(filepath.Ext(filename))])
			if err != nil {
				return nil, err
			}
			stylePath = filepath.Join(filepath.Dir(rel), formatFilename+filepath.Ext(filename))
		}

		gsf = append(gsf, &GeneratedFile{
			Path:    filepath.Join(Output, stylePath),
			Content: *bytes.NewBuffer(fileBytes),
			// Because this is a special directory for jzero
			// It is deleted to support generating all server code under the premise of only desc directory
			// Or if this file has been existed, just ignore write
			Skip: pathx.FileExists(filepath.Join(Output, "desc")) && strings.HasPrefix(rel, "desc"),
		})
	}
	return gsf, nil
}
