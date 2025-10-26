package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/property", app.createPropertyHandler)
	router.HandlerFunc(http.MethodGet, "/property/:id", app.showPropertyHandler)

	return router
}
