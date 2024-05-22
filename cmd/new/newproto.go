package new

import (
	"os"
	"path/filepath"

	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

func isNeedNewGoogleApiProto(protoDirName string) bool {
	protoDir, err := os.ReadDir(protoDirName)
	if err != nil {
		return false
	}

	for _, v := range protoDir {
		if v.IsDir() {
			continue
		}
		parser := rpcparser.NewDefaultProtoParser()
		parse, err := parser.Parse(filepath.Join(protoDirName, v.Name()))
		if err != nil {
			return false
		}
		for _, ps := range parse.Service {
			for _, rpc := range ps.RPC {
				for _, option := range rpc.Options {
					if option.Name == "(google.api.http)" {
						return true
					}
				}
			}
		}
	}
	return false
}
