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
	cfg            *config
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Flags
	cfg := new(config)
	cfg.registerFlags(nil)
	flag.Parse()

	// Dependencies Start
	slogOpts := new(slog.HandlerOptions)
	if cfg.debug {
		slogOpts.Level = slog.LevelDebug
		slogOpts.AddSource = true
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, slogOpts))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		cfg:            cfg,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
		templateCache:  templateCache,
	}
	// Dependencies End

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
