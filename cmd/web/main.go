package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP netwrok address")

	// this instantiates the new variables from the flags above
	flag.Parse()
	// User the http.NewServerMux() to initialise a new router
	mux := http.NewServeMux()

	// Create a file server to serve our static files from ./ui/static
	// The path is relative to project root
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// use mux.Handle to register the file server as a handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use http.ListenAndServ() fun to start a new server
	// take 2 params, address (:4000) and the servemux
	log.Printf("Starting server on %s 🚀", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
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
