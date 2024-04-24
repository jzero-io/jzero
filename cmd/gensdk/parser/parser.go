package parser

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
)

func Parse() {
	bytes, err := os.ReadFile(filepath.Join(".protosets", "credential.pb"))
	if err != nil {
		panic(err)
	}
	var fileSet descriptor.FileDescriptorSet
	if err := proto.Unmarshal(bytes, &fileSet); err != nil {
		return
	}
}
