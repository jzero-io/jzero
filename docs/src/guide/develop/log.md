---
title: 日志打印规范
icon: mdi:math-log
star: true
order: 4
category: 开发
tag:
  - Guide
---

```go
func (l *DownloadLogic) Download(req *types.DownloadRequest) error {
    l.Logger.Infof("download req: %v", req)

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
```

可以看到日志输入如下

```json lines
{"@timestamp":"2024-04-19T11:35:21.162+08:00","caller":"file/downloadlogic.go:33","content":"download 1.txt","level":"info","span":"0b14fa2849e40b50","trace":"a5d80df568e66150ed6d461d324f05b1"}
```

其中 logc 传入了 l.ctx 在打印后可以看到这条日志带上了 trace, 可以更好的追踪.

jzero 对于日志打印的规范:

* 请采用 logic 中 自带的 logger