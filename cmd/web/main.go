package main

import (
	"log"
	"net/http"
)

func main() {
	//initialize a new serve multiplexer & register home as '/'
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) //the {$} prevents wildcard matching
	mux.HandleFunc("GET /property", viewAllProperties)
	mux.HandleFunc("GET /property/{id}", propertyView)

	//prints log message server is starting
	log.Print("starting server on :4000")

	//start a new web server, passing in TCP addr and the servemux.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
