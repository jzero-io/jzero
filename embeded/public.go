package embeded

import (
	"embed"
	"io/fs"
)

var Web embed.FS

func RootWeb() (fs.FS, error) {
	return fs.Sub(Web, "web")
}
