package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

// Define a home handler which writes a byte slice containing
// "Hello from snippet box"
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request url exactly matches "/".
	// If it doesnt return 404
	// IMPORTANT: We return from the handler. If we did not the handler
	// would keep executing
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Initialise a slice containing the paths to 2 files. The ORDER MATTERS
	// The base template must be first
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}
	// Use the template.ParseFiles function to read the template file into a
	// template set.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Print(err.Error())
		app.serverError(w, err)
		return
	}

	// we then execute the methong on the template.
	// the last param to Execute() is the dynamic data
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		app.serverError(w, err)
	}
}

// snippetView handler
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id
	// convert it to an interger using `strconv.Atoi`, if it cant convert it
	// or its less that 1 return 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specifc snippet with ID %d", id)
}

// snippetCreate
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// r.Method to check whether the request is using POST or not
	if r.Method != http.MethodPost {
		//Send 405 status code and message
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))

		// this is a shortcut for the above
		// we are passing w to http.Error that sends the reponse for us
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("snippet create"))
}
