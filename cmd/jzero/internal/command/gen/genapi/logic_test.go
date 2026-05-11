package genapi

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestPatchLogicKeepsExistingLogicWhenRewriteHandlerFalse(t *testing.T) {
	tmpDir := withTempWorkDir(t)

	logicDir := filepath.Join(tmpDir, "internal", "logic", "user")
	if err := os.MkdirAll(logicDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	existingPath := filepath.Join(logicDir, "login.go")
	existingContent := []byte(`package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"example.com/app/internal/svc"
)

type Login struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

func NewLogin(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request, w http.ResponseWriter) *Login {
	return &Login{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
		w:      w,
	}
}

func (l *Login) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	return nil, nil
}
`)
	if err := os.WriteFile(existingPath, existingContent, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	generatedPath := filepath.Join(logicDir, "login_logic.go")
	if err := os.WriteFile(generatedPath, []byte("package user\n\ntype LoginLogic struct{}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	apiFile := filepath.Join("desc", "api", "user.api")
	ja := &JzeroApi{}
	err := ja.patchLogic(LogicFile{
		Path:           generatedPath,
		DescFilepath:   apiFile,
		Handler:        "Login",
		RewriteHandler: false,
		RequestType:    spec.DefineStruct{RawName: "LoginRequest"},
		ResponseType:   spec.DefineStruct{RawName: "LoginResponse"},
	}, map[string]*spec.ApiSpec{apiFile: {}})
	if err != nil {
		t.Fatalf("patchLogic() error = %v", err)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != string(existingContent) {
		t.Fatalf("patchLogic() should keep existing logic unchanged, got:\n%s", data)
	}
	if _, err = os.Stat(generatedPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("patchLogic() should remove generated logic when final logic exists, stat err = %v", err)
	}
}
