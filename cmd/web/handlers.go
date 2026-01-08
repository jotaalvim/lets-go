package main 

import ( 
    "fmt"
    "net/http"
    "strconv"
)


func home ( w http.ResponseWriter, _ *http.Request) {
    w.Header().Add("Server","Go")
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

