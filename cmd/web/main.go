package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/emresoysuren/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	cfg            config
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Dependencies - Form Decoder
	app.formDecoder = form.NewDecoder()

	app.sessionManager = scs.New()
	app.sessionManager.Store = mysqlstore.New(db)
	app.sessionManager.Lifetime = 12 * time.Hour
	app.sessionManager.Cookie.Secure = true

	tlsConfig := tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Serve
	srv := &http.Server{
		Addr:         app.cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		TLSConfig:    &tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", slog.String("addr", srv.Addr))

	err = srv.ListenAndServeTLS(app.cfg.tlsCert, app.cfg.tlsKey)
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
