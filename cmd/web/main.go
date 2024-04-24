package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	cfg    config
}

func main() {
	app := new(application)

	// Flags
	app.cfg.registerFlags(nil)
	flag.Parse()

	// Dependencies
	slogOpts := new(slog.HandlerOptions)
	if app.cfg.debug {
		slogOpts.Level = slog.LevelDebug
		slogOpts.AddSource = true
	}
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, slogOpts))

	// Serve
	app.logger.Info("starting server", slog.String("addr", app.cfg.addr))
	err := http.ListenAndServe(app.cfg.addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
