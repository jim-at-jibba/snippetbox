package main

import (
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containin app routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server to serve our static files from ./ui/static
	// The path is relative to project root
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// use mux.Handle to register the file server as a handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used on every request
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
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
