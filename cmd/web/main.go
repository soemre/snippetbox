package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/emresoysuren/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger        *slog.Logger
	cfg           config
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	app := new(application)

	// Flags
	app.cfg.registerFlags(nil)
	flag.Parse()

	// Dependencies - Logger
	slogOpts := new(slog.HandlerOptions)
	if app.cfg.debug {
		slogOpts.Level = slog.LevelDebug
		slogOpts.AddSource = true
	}
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, slogOpts))

	// Dependencies - DB
	db, err := openDB(app.cfg.dsn)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.templateCache = templateCache

	// Dependencies - Models
	app.snippets = &models.SnippetModel{DB: db}

	// Serve
	app.logger.Info("starting server", slog.String("addr", app.cfg.addr))
	err = http.ListenAndServe(app.cfg.addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
