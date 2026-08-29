package ui

//go:generate ../../scripts/build.sh

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:public
var publicAssets embed.FS

func Handler() http.Handler {
	f, err := fs.Sub(publicAssets, "public")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(f))
}
