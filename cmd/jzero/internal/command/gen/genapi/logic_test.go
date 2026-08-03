package genapi

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func TestPatchLogicPreservesSSESignature(t *testing.T) {
	tmpDir := withTempWorkDir(t)
	setTypesDir(t, defaultTypesDir)

	logicDir := filepath.Join(tmpDir, "internal", "logic", "conversation")
	if err := os.MkdirAll(logicDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	generatedPath := filepath.Join(logicDir, "streamevents_logic.go")
	generatedContent := []byte(`package conversation

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"example.com/app/internal/svc"
	"example.com/app/internal/types"
)

type StreamEventsLogic struct {
	logx.Logger
	ctx context.Context
	svcCtx *svc.ServiceContext
}

func NewStreamEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamEventsLogic {
	return &StreamEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx: ctx,
		svcCtx: svcCtx,
	}
}

func (l *StreamEventsLogic) StreamEvents(req *types.EventsRequest, client chan<- *types.EventsResponse) error {
	return nil
}
`)
	if err := os.WriteFile(generatedPath, generatedContent, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	apiFile := filepath.Join("desc", "api", "conversation.api")
	ja := &JzeroApi{Module: "example.com/app"}
	err := ja.patchLogic(LogicFile{
		Package:        "conversation",
		Path:           generatedPath,
		DescFilepath:   apiFile,
		Handler:        "StreamEvents",
		RewriteHandler: true,
		SSE:            true,
		RequestType:    spec.DefineStruct{RawName: "EventsRequest"},
		ResponseType:   spec.DefineStruct{RawName: "EventsResponse"},
	}, map[string]*spec.ApiSpec{apiFile: {}})
	if err != nil {
		t.Fatalf("patchLogic() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(logicDir, "streamevents.go"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	got := string(data)
	for _, want := range []string{
		"type StreamEvents struct",
		"func NewStreamEvents(",
		"client chan<- *types.EventsResponse",
		`types "example.com/app/internal/types/conversation"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("patchLogic() missing %q in SSE logic:\n%s", want, got)
		}
	}
	if strings.Contains(got, "net/http") {
		t.Fatalf("patchLogic() should not add HTTP dependencies to SSE logic:\n%s", got)
	}
}
