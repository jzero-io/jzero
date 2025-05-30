package genrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/color"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type ServerFile struct {
	DescFilepath string
	Path         string
	Service      string
}

func (jr *JzeroRpc) genServer(serverImports, pbImports ImportLines, registerServers RegisterLines) error {
	fmt.Printf("%s to generate internal/server/server.go\n", color.WithColor("Start", color.FgGreen))
	serverFile, err := templatex.ParseTemplate(map[string]any{
		"Module":          jr.Module,
		"ServerImports":   serverImports,
		"PbImports":       pbImports,
		"RegisterServers": registerServers,
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "rpc", "server.go.tpl")))
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(config.C.Wd(), "internal", "server", "server.go"), serverFile, 0o644)
	if err != nil {
		return err
	}
	fmt.Printf("%s", color.WithColor("Done\n", color.FgGreen))
	return nil
}

func (jr *JzeroRpc) GetAllServerFiles(descFilepath string, protoSpec rpcparser.Proto) ([]ServerFile, error) {
	var serverFiles []ServerFile
	for _, service := range protoSpec.Service {
		namingFormat, err := format.FileNamingFormat(config.C.Gen.Style, service.Name+"Server")
		if err != nil {
			return nil, err
		}
		fp := filepath.Join(config.C.Wd(), "internal", "server", strings.ToLower(service.Name), namingFormat+".go")
		if config.C.Gen.RpcStylePatch {
			serverDir, _ := format.FileNamingFormat(config.C.Gen.Style, service.Name)
			fp = filepath.Join(config.C.Wd(), "internal", "server", strings.ToLower(serverDir), namingFormat+".go")
		}

		f := ServerFile{
			DescFilepath: descFilepath,
			Path:         fp,
			Service:      service.Name,
		}

		serverFiles = append(serverFiles, f)
	}
	return serverFiles, nil
}
