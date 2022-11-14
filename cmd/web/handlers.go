package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	w.Write([]byte("Hello"))
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
