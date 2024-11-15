package main

import (
	"errors"
	"mango-monopoly/internal/models"
	"net/http"
	"strconv"
)

// home handler with a byte slice as the response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl", templateData{})
}

func (app *application) viewAllProperties(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, http.StatusOK, "properties.tmpl", templateData{})
}

func (app *application) propertyView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	property, err := app.properties.Get(int64(id))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.render(w, r, http.StatusOK, "property.tmpl", templateData{
		Property:      *property,
		PropertyModel: &models.PropertyModel{},
	})
}

func (app *application) createProperty(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "createProperty.tmpl", templateData{})
}
