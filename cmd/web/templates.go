package main

import "github.com/jim-at-jibba/snippetbox/internal/models"

// Define a templateData type to act as a holder for dynamic data that
// we want to pass to our templates.
// INFO: This is because `html/template` package can only accept 1 item of
// dynamic data

type templateData struct {
	Snippet *models.Snippet
}
