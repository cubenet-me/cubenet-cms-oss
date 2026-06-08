package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/htmx.min.js static/output.css
var staticFS embed.FS

func staticHandler() http.Handler {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}
