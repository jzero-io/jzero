package daemon

import (
	"fmt"
)

func Start(cfgFile string) {
	go func() {
		fmt.Println("start")
	}()
}