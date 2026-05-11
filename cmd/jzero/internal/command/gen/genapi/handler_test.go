package genapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func withTempWorkDir(t *testing.T) string {
	t.Helper()

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

	return tmpDir
}

func TestCleanHandlersDirKeepsHandlersWhenRewriteHandlerFalse(t *testing.T) {
	withTempWorkDir(t)

	handlerPath := filepath.Join("internal", "handler", "user", "custom.go")
	if err := os.MkdirAll(filepath.Dir(handlerPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(handlerPath, []byte("package user\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	apiFile := filepath.Join("desc", "api", "user.api")
	apiSpecMap := map[string]*spec.ApiSpec{
		apiFile: {
			Service: spec.Service{
				Groups: []spec.Group{
					{
						Annotation: spec.Annotation{
							Properties: map[string]string{
								"group":           "user",
								"rewrite_handler": "false",
							},
						},
					},
				},
			},
		},
	}

	ja := &JzeroApi{}
	if err := ja.cleanHandlersDir([]string{apiFile}, apiSpecMap); err != nil {
		t.Fatalf("cleanHandlersDir() error = %v", err)
	}

	if _, err := os.Stat(handlerPath); err != nil {
		t.Fatalf("cleanHandlersDir() should keep handler when rewrite_handler is false, stat err = %v", err)
	}
}

func TestPatchHandlerKeepsExistingHandlerWhenRewriteHandlerFalse(t *testing.T) {
	tmpDir := withTempWorkDir(t)

	handlerDir := filepath.Join(tmpDir, "internal", "handler", "user")
	if err := os.MkdirAll(handlerDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	existingPath := filepath.Join(handlerDir, "foo.go")
	existingContent := []byte("package user\n\nfunc Foo(  ){}\n")
	if err := os.WriteFile(existingPath, existingContent, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	generatedPath := filepath.Join(handlerDir, "foo_handler.go")
	if err := os.WriteFile(generatedPath, []byte("package user\n\nfunc FooHandler() {}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	ja := &JzeroApi{}
	err := ja.patchHandler(HandlerFile{
		Path:           generatedPath,
		Handler:        "Foo",
		RewriteHandler: false,
	}, nil)
	if err != nil {
		t.Fatalf("patchHandler() error = %v", err)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != string(existingContent) {
		t.Fatalf("patchHandler() should keep existing handler unchanged, got:\n%s", data)
	}
	if _, err = os.Stat(generatedPath); !os.IsNotExist(err) {
		t.Fatalf("patchHandler() should remove generated handler when final handler exists, stat err = %v", err)
	}
}
