package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Display a snippet"))
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

func main() {
	// User the http.NewServerMux() to initialise a new router
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use http.ListenAndServ() fun to start a new server
	// take 2 params, address (:4000) and the servemux
	log.Print("Starting server on :4000 🚀")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
