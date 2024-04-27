package main

import "flag"

type config struct {
	addr      string
	dsn       string
	staticDir string
	debug     bool
	tlsCert   string
	tlsKey    string
}

// If flagSet is nil, it will be used as flag.CommandLine by default.
func (cfg *config) registerFlags(flagSet *flag.FlagSet) {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}
	flagSet.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flagSet.StringVar(&cfg.dsn, "dsn", "web:admin@/snippetbox?parseTime=true", "MySQL data source name")
	flagSet.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flagSet.BoolVar(&cfg.debug, "debug", false, "Debug Mode")
	flagSet.StringVar(&cfg.tlsCert, "tls-cert", "./tls/cert.pem", "Path to TLS Certificate")
	flagSet.StringVar(&cfg.tlsKey, "tls-key", "./tls/key.pem", "Path to TLS Key")
}
