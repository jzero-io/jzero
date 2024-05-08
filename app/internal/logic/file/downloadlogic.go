package file

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/app/internal/svc"
	"github.com/jzero-io/jzero/app/internal/types"
)

type DownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer io.Writer
}

func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer io.Writer) *DownloadLogic {
	return &DownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}
}

func (l *DownloadLogic) Download(req *types.DownloadRequest) error {
	l.Logger.Infof("download req: %s", req.File)
	l.Logger.Error("test fail")

	body, err := os.ReadFile(filepath.Join("./filedata", req.File))
	if err != nil {
		return err
	}

	n, err := l.writer.Write(body)
	if err != nil {
		return err
	}

	if n < len(body) {
		return io.ErrClosedPipe
	}

	return nil
}
