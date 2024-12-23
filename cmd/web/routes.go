package main

import (
	"net/http"

	"github.com/emresoysuren/snippetbox/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initial
	mux := http.NewServeMux()

	// File Servers
	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /ping", ping)

	// Handlers - Unprotected
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Handlers - Protected
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standart := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standart.Then(mux)
}
