package genapi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestSeparateTypesGoRewritesDefaultTypesFileWhenNoDefaultTypesRemain(t *testing.T) {
	tmpDir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}

	if err := os.MkdirAll(filepath.Join("internal", "types"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	oldContent := []byte(`package types

type OldType struct{}
`)
	defaultTypesPath := filepath.Join("internal", "types", "types.go")
	if err := os.WriteFile(defaultTypesPath, oldContent, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	apiSpecMap := map[string]*spec.ApiSpec{
		"desc/api/user.api": {
			Info: spec.Info{
				Properties: map[string]string{
					"go_package": "user",
				},
			},
			Types: []spec.Type{
				spec.DefineStruct{
					RawName: "UserReq",
					Members: []spec.Member{
						{
							Name: "Name",
							Type: spec.PrimitiveType{RawName: "string"},
							Tag:  "`json:\"name\"`",
						},
					},
				},
			},
		},
	}

	ja := &JzeroApi{}
	if err := ja.separateTypesGo([]string{"desc/api/user.api"}, apiSpecMap); err != nil {
		t.Fatalf("separateTypesGo() error = %v", err)
	}

	data, err := os.ReadFile(defaultTypesPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	got := string(data)
	if strings.Contains(got, "OldType") {
		t.Fatalf("separateTypesGo() should rewrite stale default types file, got:\n%s", got)
	}
	if !strings.Contains(got, "package types") {
		t.Fatalf("separateTypesGo() should keep default types package, got:\n%s", got)
	}
}
