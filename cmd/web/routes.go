package main

import "net/http"

func (app *application) routes() *http.ServeMux {
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

	return mux
}