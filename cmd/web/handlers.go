package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// Define a home handler which writes a byte slice containing
// "Hello from snippet box"
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request url exactly matches "/".
	// If it doesnt return 404
	// IMPORTANT: We return from the handler. If we did not the handler
	// would keep executing
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialise a slice containing the paths to 2 files. The ORDER MATTERS
	// The base template must be first
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}
	// Use the template.ParseFiles function to read the template file into a
	// template set.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// we then execute the methong on the template.
	// the last param to Execute() is the dynamic data
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}

// snippetView handler
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id
	// convert it to an interger using `strconv.Atoi`, if it cant convert it
	// or its less that 1 return 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specifc snippet with ID %d", id)
}

// snippetCreate
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// r.Method to check whether the request is using POST or not
	if r.Method != http.MethodPost {
		//Send 405 status code and message
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))

		// this is a shortcut for the above
		// we are passing w to http.Error that sends the reponse for us
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("snippet create"))
}
