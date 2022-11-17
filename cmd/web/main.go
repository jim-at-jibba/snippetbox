package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Define an application struct to hold the application-wide dependencies for
// the web app.
type applicaion struct {
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
	app := &applicaion{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// User the http.NewServerMux() to initialise a new router
	mux := http.NewServeMux()

	// Create a file server to serve our static files from ./ui/static
	// The path is relative to project root
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// use mux.Handle to register the file server as a handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	// Use http.ListenAndServ() fun to start a new server
	// take 2 params, address (:4000) and the servemux
	infoLog.Printf("Starting server on %s 🚀", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
