package genrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/color"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/embeded"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type ServerFile struct {
	Path    string
	Service string
}

func (jr *JzeroRpc) genServer(serverImports, pbImports jzerodesc.ImportLines, registerServers jzerodesc.RegisterLines) error {
	fmt.Printf("%s to generate internal/server/server.go\n", color.WithColor("Start", color.FgGreen))
	serverFile, err := templatex.ParseTemplate(map[string]any{
		"Module":          jr.Module,
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("app", "internal", "server", "server.go.tpl")))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(jr.Wd, "internal", "server", "server.go"), serverFile, 0o644)
	if err != nil {
		return err
	}
	fmt.Printf("%s", color.WithColor("Done\n", color.FgGreen))
	return nil
}

func (jr *JzeroRpc) GetAllServerFiles(protoSpec rpcparser.Proto) ([]ServerFile, error) {
	var serverFiles []ServerFile
	for _, service := range protoSpec.Service {
		namingFormat, err := format.FileNamingFormat(jr.Style, service.Name+"Server")
		if err != nil {
			return nil, err
		}
		fp := filepath.Join(jr.Wd, "internal", "server", strings.ToLower(service.Name), namingFormat+".go")
		if jr.RpcStylePatch {
			serverDir, _ := format.FileNamingFormat(jr.Style, service.Name)
			fp = filepath.Join(jr.Wd, "internal", "server", strings.ToLower(serverDir), namingFormat+".go")
		}

		f := ServerFile{
			Path:    fp,
			Service: service.Name,
		}

		serverFiles = append(serverFiles, f)
	}
	return serverFiles, nil
}
