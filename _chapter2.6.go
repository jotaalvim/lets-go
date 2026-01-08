package main 

import ( 
    //"reflect"
    "fmt"
    "log"
    "net/http"
    "strconv"
)


// Define a home handler function which writes a byte slice containing
func home ( w http.ResponseWriter, _ *http.Request) {
    w.Header().Add("Server","Go")
    w.Header().Set("Cache-Control","public,max-age=3456")
    // w.Header().Get("Cache-Control")
    // w.Header().Valyes("Cache-Control") // devolve uma slice

    w.Write( [] byte ("Hello world"))
}

func create ( w http.ResponseWriter, _ *http.Request) {
    w.Write ( [] byte ("create something"))
}

func createPost ( w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader (http.StatusCreated) // cria status code sozinho
    w.Write([] byte ("save a new snippet"))
}

func view ( w http.ResponseWriter, r *http.Request) {
    id,err := strconv.Atoi( r.PathValue("id") ) 
    if err != nil || id < 1 {
        http.NotFound(w,r)
        return 
    }   
    message := fmt.Sprintf("display something =  %d", id)
    w.Write ( [] byte (message) )
}



func main () {

    handler := http.NewServeMux()


    handler.HandleFunc("GET  /{$}"      , home       ) // Restrict this route to exact matches on / only.
    handler.HandleFunc("GET  /view/{id}", view       )
    handler.HandleFunc("GET  /create"   , create     )
    handler.HandleFunc("POST /create"   , createPost )

    log.Println(" starting server on port 4000", handler)

    //fmt.Println(reflect.TypeOf(*handler))

    err := http.ListenAndServe( ":4000", handler)

    log.Fatal(err)


}
