package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/emresoysuren/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Paths Start
	var (
		htmlDirPath  = filepath.Join(".", "ui", "html")
		basePath     = filepath.Join(htmlDirPath, "base.html")
		pagesPath    = filepath.Join(htmlDirPath, "pages", "*.html")
		partialsPath = filepath.Join(htmlDirPath, "partials", "*.html")
	)

	pages, err := filepath.Glob(pagesPath)
	if err != nil {
		return nil, err
	}
	// Paths End

	// Template Set
	ts, err := template.New("base").Funcs(functions).ParseFiles(basePath)
	if err != nil {
		return nil, err
	}
	ts, err = ts.ParseGlob(partialsPath)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		pagets, err := ts.Clone()
		if err != nil {
			return nil, err
		}

		pagets, err = pagets.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = pagets
	}

	return cache, nil
}
