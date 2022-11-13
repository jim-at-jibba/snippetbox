package main

import (
	"log"
	"net/http"
)

// Define a home handler which writes a byte slice containing
// "Hello from snippet box"
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func main() {
	// User the http.NewServerMux() to initialise a new router
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Use http.ListenAndServ() fun to start a new server
	// take 2 params, address (:4000) and the servemux
	log.Print("Starting server on :4000 🚀")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
