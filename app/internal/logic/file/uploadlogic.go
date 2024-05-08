package file

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/app/internal/svc"
	"github.com/jzero-io/jzero/app/internal/types"
)

const maxFileSize = 10 << 20 // 10 MB

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload() (resp *types.UploadResponse, err error) {
	err = l.r.ParseMultipartForm(maxFileSize)
	if err != nil {
		return nil, err
	}
	file, handler, err := l.r.FormFile("myFile")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logx.Infof("upload file: %+v, file size: %d, MIME header: %+v",
		handler.Filename, handler.Size, handler.Header)

	tempFile, err := os.Create(path.Join("./filedata", handler.Filename))
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, err
	}

	return &types.UploadResponse{
		Code: 0,
	}, nil
}
