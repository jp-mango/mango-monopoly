package main

import (
	"net/http"

	chain "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//initialize a new serve multiplexer & register home as '/'
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home) //the {$} prevents wildcard matching

	mux.HandleFunc("GET /properties", app.viewAllProperties)

	mux.HandleFunc("GET /property/{id}", app.propertyView)

	mux.HandleFunc("GET /property/create", app.createPropertyPage)
	mux.HandleFunc("POST /property/create", app.propertyCreatePost)

	standard := chain.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
