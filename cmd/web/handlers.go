package main

import (
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
