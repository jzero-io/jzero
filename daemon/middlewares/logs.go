package middlewares

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/nxadm/tail"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/jaronnie/jzero/daemon/internal/config"
)

func PrintLogToConsole(c config.Config) {
	if c.Log.Mode == "console" {
		return
	}

	if !c.Jzero.LogToConsole {
		return
	}

	logs := []string{"access.log", "error.log", "severe.log", "slow.log"}
	if c.Log.Stat {
		logs = append(logs, "stat.log")
	}

	var logPaths []string
	for _, v := range logs {
		logPaths = append(logPaths, filepath.Join(c.Log.Path, v))
	}

	go func() {
		err := mr.MapReduceVoid(func(source chan<- string) {
			for _, v := range logPaths {
				source <- v
			}
		}, func(item string, writer mr.Writer[*tail.Tail], cancel func(error)) {
			t, err := tail.TailFile(
				filepath.Join(item), tail.Config{
					Follow: true,
					ReOpen: true,
					Poll:   true,
					Location: &tail.SeekInfo{
						Offset: 0,
						Whence: io.SeekEnd,
					},
					Logger: tail.DiscardingLogger,
				})
			if err != nil {
				cancel(err)
			}

			// Print the text of each received line
			for line := range t.Lines {
				fmt.Println(line.Text)
			}

		}, func(pipe <-chan *tail.Tail, cancel func(error)) {}, mr.WithWorkers(len(logs)))

		if err != nil {
			return
		}
	}()
}
