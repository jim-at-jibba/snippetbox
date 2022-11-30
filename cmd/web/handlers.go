package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jim-at-jibba/snippetbox/internal/models"
	"github.com/jim-at-jibba/snippetbox/internal/validator"
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

	// giving the expires a default value
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// Embedding the Validator means the snippetCreateForm "inherits" all the
// fields and methods of our Validator "Composition over inheritance"
// struct tags tell the decoder how to map HTML form value to the struct fields
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // this means ignore this
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

	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// Becaise values come from .Get as a string and we need a number we need to
	// convert it
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// replaced with decoder
	// form := snippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// }
	var form snippetCreateForm

	err := app.decodeePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// fieldErrors := make(map[string]string)
	// if strings.TrimSpace(form.Title) == "" {
	// 	form.FieldErrors["title"] = "This field cannot be blank"
	// 	// utf8.Rune* counts the characters. Go len would count the bytes
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	form.FieldErrors["title"] = "This field cannot be over 100 characters long"
	// }
	//
	// if strings.TrimSpace(form.Content) == "" {
	// 	form.FieldErrors["content"] = "This field cannot be blank"
	// }
	//
	// if expires != 1 && expires != 7 && expires != 365 {
	// 	form.FieldErrors["expires"] = "This field must equal, 1, 7 or 365"
	// }
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 chars long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 ,365")

	// Use the Valid() method to see if there are any messages, re-display the create.tmpl.html,
	// passing the snippetCreatePost instance as dynic data
	// and return from the handler
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
