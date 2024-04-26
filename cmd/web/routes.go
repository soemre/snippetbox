package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initial
	mux := http.NewServeMux()

	// File Servers
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(app.cfg.staticDir)})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Handlers
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	standart := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standart.Then(mux)
}
