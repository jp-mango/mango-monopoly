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

	app.render(w, r, http.StatusOK, "property.tmpl", templateData{
		Property:      *property,
		PropertyModel: &models.PropertyModel{},
	})
}

// TODO: accept user input
func (app *application) createProperty(w http.ResponseWriter, r *http.Request) {
	//dummy data
	property := &models.Property{
		Address:         sql.NullString{String: "45 Wallaby Way."},
		City:            sql.NullString{String: "Sydney"},
		State:           sql.NullString{String: "Australia"},
		Zip:             sql.NullString{String: "1337"},
		ParcelID:        sql.NullString{String: "DU"},
		PropertyType:    sql.NullString{String: "land"},
		LandValue:       sql.NullInt64{Int64: 0},
		BuildingValue:   sql.NullInt64{Int64: 24500000},
		FairMarketValue: sql.NullInt64{Int64: 27650000},
		LotSize:         sql.NullFloat64{Float64: 102.8},
	}

	id, err := app.properties.Insert(property)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/property/%d", id), http.StatusSeeOther)
}
