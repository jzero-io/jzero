package fuzzy

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func DecodeRequest(r *http.Request, req any) error {
	if r.Body == nil {
		return nil
	}

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		return nil
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err = r.Body.Close(); err != nil {
		return err
	}

	bodyBytes, err = Decode(bodyBytes, req)
	if err != nil {
		return err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	r.ContentLength = int64(len(bodyBytes))

	logx.Debugf("new request body bytes: %s", bodyBytes)

	return nil
}
