package main

import (
    "net/http"
)

func (app *application) routes () *http.ServeMux {

    handler    := http.NewServeMux()

    fileServer := http.FileServer( http.Dir( app.cfg.staticDir ))

    handler.Handle(    "GET /static/"   , http.StripPrefix("/static", fileServer)   )

    handler.HandleFunc("GET  /{$}"      , app.home       ) 
    handler.HandleFunc("GET  /view/{id}", app.view       )
    handler.HandleFunc("GET  /create"   , app.create     )
    handler.HandleFunc("POST /create"   , app.createPost )

    return handler

}
