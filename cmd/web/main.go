package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for
// the web app.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP netwrok address")

	// this instantiates the new variables from the flags above
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialise a new instance of our applicaion struct, containing the deps
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// Use http.ListenAndServ() fun to start a new server
	// take 2 params, address (:4000) and the servemux
	infoLog.Printf("Starting server on %s 🚀", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
