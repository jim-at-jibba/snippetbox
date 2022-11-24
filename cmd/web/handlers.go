package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jim-at-jibba/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

// Define a home handler which writes a byte slice containing
// "Hello from snippet box"
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// The following was only needed when we were not using a full
	// featured router - leaving for future knowledge
	// Check if the current request url exactly matches "/".
	// If it doesnt return 404
	// IMPORTANT: We return from the handler. If we did not the handler
	// would keep executing
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Initialise a slice containing the paths to 2 files. The ORDER MATTERS
	// The base template must be first
	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }
	// // Use the template.ParseFiles function to read the template file into a
	// // template set.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }
	//
	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// we then execute the methong on the template.
	// the last param to Execute() is the dynamic data
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// }

	// Call the newTemplateData() helper to get a templateData struct
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

// snippetView handler
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	// Extract the value of the id
	// convert it to an interger using `strconv.Atoi`, if it cant convert it
	// or its less that 1 return 404
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		// Current best practice for checking for specific error types
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// snippetCreate
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Not needed now we have a full featured router
	// // r.Method to check whether the request is using POST or not
	// if r.Method != http.MethodPost {
	// 	//Send 405 status code and message
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	// w.WriteHeader(405)
	// 	// w.Write([]byte("Method not allowed"))
	//
	// 	// this is a shortcut for the above
	// 	// we are passing w to http.Error that sends the reponse for us
	// 	// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	title := "O snail"
	content := "O snial, where are these snippets coming from?"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
