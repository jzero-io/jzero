package new

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

var (
	base64Matcher = regexp.MustCompile(`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`)
)

func IsBase64(base64 string) bool {
	return base64Matcher.MatchString(base64)
}

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

func Run(appName, base string) error {
	if err := os.MkdirAll(config.C.New.Output, 0o755); err != nil {
		return err
	}

	templateData, err := NewTemplateData()
	if err != nil {
		return err
	}

	templateData["Features"] = config.C.New.Features
	templateData["Module"] = config.C.New.Module
	templateData["APP"] = appName
	if abs, err := filepath.Abs(config.C.New.Output); err == nil {
		templateData["DirName"] = filepath.Base(abs)
	} else {
		return err
	}
	templateData["Style"] = config.C.New.Style

	jn := JzeroNew{
		TemplateData: templateData,
		nc:           config.C.New,
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
		bytesFormat, err = gosimports.Process("", bytes, nil)
		if err != nil {
			return errors.Wrapf(err, "format %s", path)
		}
	}

	// 增加可执行权限
	if lo.Contains(config.C.New.ExecutableExtensions, filepath.Ext(path)) {
		return os.WriteFile(path, bytesFormat, 0o744)
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

		filename := file.Name()
		if IsBase64(filename) {
			filenameBytes, _ := base64.StdEncoding.DecodeString(filename)
			filename = string(filenameBytes)
		}

		filename = strings.TrimSuffix(filename, ".tpl")

		rel, err := filepath.Rel(jn.base, filepath.Join(dirname, filename))
		if err != nil {
			return nil, err
		}

		var fileBytes []byte
		if strings.HasSuffix(file.Name(), ".tpl.tpl") {
			// .tpl.tpl suffix means it is a template, do not parse if anymore
			fileBytes = embeded.ReadTemplateFile(filepath.Join(dirname, file.Name()))
		} else {
			fileBytes, err = templatex.ParseTemplate(filepath.Join(dirname, file.Name()), jn.TemplateData, embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())))
			if err != nil {
				return nil, err
			}
		}

		// parse template name
		templatePath := filepath.Join(filepath.Dir(rel), filename)
		stylePathBytes, err := templatex.ParseTemplate(templatePath, jn.TemplateData, []byte(templatePath))
		if err != nil {
			return nil, err
		}

		// specify
		if filename == "go.mod" && jn.nc.Mono {
			continue
		}

		gsf = append(gsf, &GeneratedFile{
			Path:    filepath.Join(jn.nc.Output, string(stylePathBytes)),
			Content: *bytes.NewBuffer(fileBytes),
			Skip: func() bool {
				var ignore []string
				for _, v := range jn.nc.Ignore {
					ignore = append(ignore, filepath.ToSlash(v))
				}
				for _, v := range jn.nc.IgnoreExtra {
					ignore = append(ignore, filepath.ToSlash(v))
				}
				return lo.Contains(ignore, rel)
			}(),
		})
	}
	return gsf, nil
}
