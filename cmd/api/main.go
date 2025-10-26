package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		//AddSource: true,
	}))

	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug("No .env file found; relying on existing environment variables")
		} else {
			logger.Warn("Error loading environment variables", "error", err.Error())
		}
	}

	var cfg config

	srvPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Error("unable to convert port env")
	}
	flag.IntVar(&cfg.port, "port", srvPort, "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "Environment (development|staging|production)")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
