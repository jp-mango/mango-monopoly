package main

import (
	"fmt"
	"html/template"
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

	files := []string{
		"html/base.tmpl",
		"html/pages/property.tmpl",
		"html/partials/nav.tmpl",
	}

	ts, err := template.ParseFS(ui.TemplateFiles, files...)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", id)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// TODO: fix handler - curl request returning 0 for id
func (app *application) createProperty(w http.ResponseWriter, r *http.Request) {
	//dummy data
	addr := "123 test rd."
	city := "townsville"
	state := "Georgia"
	zip := "1234"
	parselID := "abcd"
	propertyType := "Multifamily compound"
	landValue := 58394054
	buildingValue := 54355324
	fmv := 4325454254234
	lotSize := 43.6

	id, err := app.properties.Insert(addr, city, state, zip, parselID, propertyType, float32(landValue), float32(buildingValue), float32(fmv), float32(lotSize))
	if err != nil {
		app.serverError(w, r, err)
		fmt.Println("Error inserting property:", err)
		return
	}

	if id == 0 {
		fmt.Println("Insert successful but ID is 0, possible issue with query or RETURNING clause.")
	}

	http.Redirect(w, r, fmt.Sprintf("/property/%d", id), http.StatusSeeOther)
}
