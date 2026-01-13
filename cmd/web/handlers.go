package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"modulo.porreiro/internal/models"
)

///foo/bar?title=value&content=value .
//You can retrieve the values for the query string parameters in your handlers via the
//r.URL.Query().Get()

// posso aceder ao app logger dentro da função
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//panic("panico oops aqui vau bomba")
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)

}

func (app *application) create(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	// Limit the request body size to 4096 bytes
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	// For all requests, ParseForm parses the raw query from the URL and updates r.Form
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newPath := fmt.Sprintf("/view/%d", id)
	http.Redirect(w, r, newPath, http.StatusSeeOther)
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}
