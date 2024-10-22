package main

import (
	"html/template"
	"log"
	"mango-monopoly/ui"
	"net/http"
	"strconv"
)

// home handler with a byte slice as the response body
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"html/base.tmpl",
		"html/pages/home.tmpl",
		"html/partials/nav.tmpl",
	}

	ts, err := template.ParseFS(ui.TemplateFiles, files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func viewAllProperties(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"html/base.tmpl",
		"html/pages/properties.tmpl",
		"html/partials/nav.tmpl",
	}

	ts, err := template.ParseFS(ui.TemplateFiles, files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func propertyView(w http.ResponseWriter, r *http.Request) {
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
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
