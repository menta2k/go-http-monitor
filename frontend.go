package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func frontendHandler() http.Handler {
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		f, err := distFS.Open(r.URL.Path[1:]) // strip leading /
		if err != nil {
			// SPA fallback: serve index.html for any non-file route
			r.URL.Path = "/"
		} else {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}
