package main

import (
	"html/template"
	"path/filepath"

	"github.com/jim-at-jibba/snippetbox/internal/models"
)

// Define a templateData type to act as a holder for dynamic data that
// we want to pass to our templates.
// INFO: This is because `html/template` package can only accept 1 item of
// dynamic data

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialise a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths
	// that match the pattern "./ui/html/pages/*.tmpl". This will eventually
	// gives us a slice of all the filepaths for our application page templates
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name and assign it to the name variable
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		cache[name] = ts
	}

	return cache, nil
}
