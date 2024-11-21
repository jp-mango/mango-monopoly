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

	dynamic := chain.New(app.sessionManager.LoadAndSave, noSurf)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) //the {$} prevents wildcard matching
	mux.Handle("GET /property/{id}", dynamic.ThenFunc(app.propertyView))
	mux.Handle("GET /properties", dynamic.ThenFunc(app.viewAllProperties))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /property/create", protected.ThenFunc(app.createPropertyPage))
	mux.Handle("POST /property/create", protected.ThenFunc(app.propertyCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := chain.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
