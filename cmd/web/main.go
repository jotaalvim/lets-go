package main 

import ( 
    "log"
    "net/http"
)


func main () {

    handler := http.NewServeMux()

    
    fileServer := http.FileServer( http.Dir("./ui/static/"))

    handler.Handle(    "GET /static/"   , http.StripPrefix("/static", fileServer)   )

    handler.HandleFunc("GET  /{$}"      , home       ) 
    handler.HandleFunc("GET  /view/{id}", view       )
    handler.HandleFunc("GET  /create"   , create     )
    handler.HandleFunc("POST /create"   , createPost )



    log.Println("starting server on port 4000", handler)
    log.Println("https:://localhost:4000")

    err := http.ListenAndServe( ":4000", handler)
    //func ListenAndServe(addr string, handler Handler) error

    log.Fatal(err)
}
