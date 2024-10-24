package main

import (
	"html/template"
	"mango-monopoly/internal/models"
	"path/filepath"
)

type templateData struct {
	Property      models.Property
	PropertyModel *models.PropertyModel
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl".
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// Parse the files into a template set.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}
	return cache, nil
}
