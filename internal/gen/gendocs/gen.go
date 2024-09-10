package gendocs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/internal/gen/gendocs/markdown"
	gendocsparser "github.com/jzero-io/jzero/internal/gen/gendocs/parser"
)

func Gen(gc config.GenConfig) error {
	if pathx.FileExists(gc.Docs.ApiDir) {
		mainApiFile, isDelete, err := gen.GetMainApiFilePath(filepath.Join("desc", "api"))
		if err != nil {
			return err
		}
		defer func() {
			if isDelete {
				_ = os.Remove(mainApiFile)
			}
		}()

		p, err := parser.Parse(mainApiFile, nil)
		if err != nil {
			return err
		}

		var docsSpecs []*gendocsparser.DocsSpec

		var groups []string
		for _, v := range p.Service.Groups {
			groups = append(groups, v.GetAnnotation("group"))
		}

		docsParser := gendocsparser.NewDocsParser(p)

		docsSpecs = docsParser.BuildDocsSpecHierarchy(groups)

		m := markdown.New(docsSpecs)
		err = m.Generate()
		if err != nil {
			return err
		}
	}

	if pathx.FileExists(gc.Docs.ProtoDir) {
		_ = os.MkdirAll(gc.Docs.Output, 0o755)
		protoFilepath, err := gen.GetProtoFilepath(gc.Swagger.ProtoDir)
		if err != nil {
			return err
		}

		command := fmt.Sprintf("protoc -I%s -I%s --doc_out=%s --doc_opt=%s,index.%s %s",
			gc.Docs.ProtoDir,
			filepath.Join(gc.Docs.ProtoDir, "third_party"),
			gc.Docs.Output,
			gc.Docs.Format,
			getExt(gc.Docs.Format),
			strings.Join(protoFilepath, " "),
		)
		_, err = execx.Run(command, gc.Swagger.Wd())
		if err != nil {
			return err
		}
	}

	return nil
}

func getExt(format string) string {
	switch format {
	case "markdown":
		return "md"
	}
	return format
}
