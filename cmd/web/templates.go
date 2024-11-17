package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"mango-monopoly/internal/models"
	"mango-monopoly/ui"
	"path/filepath"
)

type templateData struct {
	CurrentYear   int
	Property      models.Property
	Properties    []models.Property
	PropertyModel *models.PropertyModel
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the embed.FS to get a slice of all filepaths that match the pattern.
	// Check if 'fs.Glob()' returns any errors or no files.
	pages, err := fs.Glob(ui.TemplateFiles, "html/pages/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("error finding pages: %w", err)
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no page templates found in embedded files")
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepaths for our base template, any
		// partials and the page itself.
		files := []string{
			"html/base.tmpl",
			"html/partials/nav.tmpl",
			page,
		}

		// Parse the files using the embedded filesystem.
		ts, err := template.ParseFS(ui.TemplateFiles, files...)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %w", name, err)
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}

	return cache, nil
}
