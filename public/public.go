package public

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var assets embed.FS

func RootAssets() (fs.FS, error) {
	return fs.Sub(assets, "dist")
}
