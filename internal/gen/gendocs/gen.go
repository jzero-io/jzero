package gendocs

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/internal/gen/gendocs/markdown"

	"github.com/jzero-io/jzero/internal/gen"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"

	gendocsparser "github.com/jzero-io/jzero/internal/gen/gendocs/parser"
)

func Gen(cmd *cobra.Command, args []string) error {
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

	return nil
}
