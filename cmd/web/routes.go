package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.cfg.staticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET  /{$}", app.home)
	mux.HandleFunc("GET  /view/{id}", app.view)
	mux.HandleFunc("GET  /create", app.create)
	mux.HandleFunc("POST /create", app.createPost)
	// onde ponho o meu 404 render
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))

}
