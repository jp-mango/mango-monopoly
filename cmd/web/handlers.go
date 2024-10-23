package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"mango-monopoly/internal/models"
	"mango-monopoly/ui"
	"net/http"
	"strconv"
)

// home handler with a byte slice as the response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"html/base.tmpl",
		"html/pages/home.tmpl",
		"html/partials/nav.tmpl",
	}

	ts, err := template.ParseFS(ui.TemplateFiles, files...)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) viewAllProperties(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"html/base.tmpl",
		"html/pages/properties.tmpl",
		"html/partials/nav.tmpl",
	}

	ts, err := template.ParseFS(ui.TemplateFiles, files...)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

	files := []string{
		"html/base.tmpl",
		"html/pages/property.tmpl",
		"html/partials/nav.tmpl",
	}

	ts, err := template.ParseFS(ui.TemplateFiles, files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Property: *property,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// TODO: accept user input
func (app *application) createProperty(w http.ResponseWriter, r *http.Request) {
	//dummy data
	property := &models.Property{
		Address:         sql.NullString{String: "45 Wallaby Way.", Valid: true},
		City:            sql.NullString{String: "Sydney", Valid: true},
		State:           sql.NullString{String: "Australia", Valid: true},
		Zip:             sql.NullString{String: "1337", Valid: true},
		ParcelID:        sql.NullString{String: "DU", Valid: true},
		PropertyType:    sql.NullString{String: "land", Valid: true},
		LandValue:       sql.NullFloat64{Float64: 0, Valid: true},
		BuildingValue:   sql.NullFloat64{Float64: 24500000, Valid: true},
		FairMarketValue: sql.NullFloat64{Float64: 27650000, Valid: true},
		LotSize:         sql.NullFloat64{Float64: 102.8, Valid: true},
	}

	id, err := app.properties.Insert(property)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/property/%d", id), http.StatusSeeOther)
}
