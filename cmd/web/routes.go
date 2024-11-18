package main

import (
	"net/http"

	"github.com/justinas/alice"
	chain "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//initialize a new serve multiplexer & register home as '/'
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) //the {$} prevents wildcard matching

	mux.Handle("GET /properties", dynamic.ThenFunc(app.viewAllProperties))

	mux.Handle("GET /property/{id}", dynamic.ThenFunc(app.propertyView))

	mux.Handle("GET /property/create", dynamic.ThenFunc(app.createPropertyPage))
	mux.Handle("POST /property/create", dynamic.ThenFunc(app.propertyCreatePost))

	standard := chain.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
