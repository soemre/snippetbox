package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		err.Error(),
		slog.String("method", r.Method),
		slog.String("uri", r.RequestURI),
		slog.String("trace", string(debug.Stack())),
	)
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	_, err := buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}
