package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"mango-monopoly/internal/models"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type application struct {
	logger        *slog.Logger
	properties    *models.PropertyModel
	templateCache map[string]*template.Template
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("DSN"), "Postgres data source name")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		//AddSource: true,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("db connection established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		properties:    &models.PropertyModel{DB: db},
		templateCache: templateCache,
	}
	//prints log message server is starting
	logger.Info("starting server", "addr", *addr)

	//start a new web server, passing in TCP addr and the servemux.
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
