package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.cfg.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// LoadAndSave provides middleware which automatically loads and saves
	// session data for the current request, and communicates the session token
	// to and from the client in a cookie.
	mux.Handle("GET  /{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("GET  /view/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.view)))
	mux.Handle("GET  /create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.create)))
	mux.Handle("POST /create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.createPost)))

	// onde ponho o meu 404 render?
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))

}
