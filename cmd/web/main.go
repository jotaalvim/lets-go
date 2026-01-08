package main 

import ( 
    "log"
    "net/http"
)


func main () {

    handler := http.NewServeMux()

    handler.HandleFunc("GET  /{$}"      , home       ) 
    handler.HandleFunc("GET  /view/{id}", view       )
    handler.HandleFunc("GET  /create"   , create     )
    handler.HandleFunc("POST /create"   , createPost )

    log.Println(" starting server on port 4000", handler)

    err := http.ListenAndServe( ":4000", handler)

    log.Fatal(err)
}
