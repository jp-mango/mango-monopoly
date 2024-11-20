package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"mango-monopoly/internal/models"
	"mango-monopoly/ui"
	"path/filepath"
	"strconv"
)

type templateData struct {
	CurrentYear     int
	Property        models.Property
	Properties      []models.Property
	Form            any
	Flash           string
	IsAuthenticated bool
}

func formatMoney(price int64) string {
	m := strconv.FormatInt(price, 10)
	n := len(m)
	if n <= 3 {
		return m
	}

	var result string
	for i := n - 1; i >= 0; i-- {
		result = string(m[i]) + result
		if (n-i)%3 == 0 && i != 0 {
			result = "," + result
		}
	}
	return result
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	funcs := template.FuncMap{
		"FormatMoney": formatMoney,
	}

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
		ts, err := template.New(name).Funcs(funcs).ParseFS(ui.TemplateFiles, files...)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %w", name, err)
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}

	return cache, nil
}
