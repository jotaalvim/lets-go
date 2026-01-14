package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"modulo.porreiro/internal/models"
	"modulo.porreiro/internal/validator"
)

type snippetCreateForm struct {
	// this tags tell the Decoder how to map the html form values into different fields
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// embeding the validator means that snippetCreateForm "inherits" all the felds and methods
	// posso aceder com form.Valid(), ou também form.Validator.Valid()
	validator.Validator `form:-`
}

///foo/bar?title=value&content=value .
//You can retrieve the values for the query string parameters in your handlers via the
//r.URL.Query().Get()

func (app *application) home(w http.ResponseWriter, r *http.Request) {

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
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	//expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	//if err != nil {
	//	app.serverError(w, r, err)
	//	return
	//}
	//form := snippetCreateForm{
	//	Title:   r.PostForm.Get("title"),
	//	Content: r.PostForm.Get("content"),
	//	Expires: expires,
	//	//FieldErrors: map[string]string{},
	//}
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 chars")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be longer than 100 chars")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		// HTTP status code 422 Unprocessable Entity
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

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

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "user signup")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "user signup post")
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "user login ")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "user login post")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "user logout post")
}
