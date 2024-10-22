package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	//initialize a new serve multiplexer & register home as '/'
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) //the {$} prevents wildcard matching
	mux.HandleFunc("GET /properties", viewAllProperties)
	mux.HandleFunc("GET /property/{id}", propertyView)

	//prints log message server is starting
	logger.Info("starting server", "addr", *addr)

	//start a new web server, passing in TCP addr and the servemux.
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
