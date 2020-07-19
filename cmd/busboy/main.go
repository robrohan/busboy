package main

import (
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

var build = "develop"

type Config struct {
	APIHost         string        `conf:"default:0.0.0.0:8580"`
	Root            string        `config:"."`
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:5s"`
	ShutdownTimeout time.Duration `conf:"default:5s"`
}

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "BB : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	cfg := Config{}

	if err := conf.Parse(os.Args[1:], "BB", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("BB", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting
	expvar.NewString("build").Set(build)
	log.Printf("Started : BusBoy initializing : version %q", build)
	defer log.Println("Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("Config :\n%v\n", out)

	// =========================================================================
	api := http.Server{
		Addr:         cfg.APIHost,
		Handler:      newAPI(log, cfg),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	log.Printf("API listening on %s", api.Addr)
	err = api.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func newAPI(log *log.Logger, cfg Config) http.Handler {
	fs := http.FileServer(http.Dir(cfg.Root))
	http.Handle("/", fs)
	return LoggingMiddleware(log)(fs)
}

func LoggingMiddleware(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.RemoteAddr, r.RequestURI)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
