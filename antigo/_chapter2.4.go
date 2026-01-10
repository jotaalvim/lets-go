package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
func home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello world"))
}

func create(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("create something"))
}

func view(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// a função string() recebe bytes e devolve uma string
	//v       := "display someting : " + string(97)
	//fmt.Println(v , id)
	//w.Write ( [] byte (v) )
	message := fmt.Sprintf("display something =  %d", id)
	w.Write([]byte(message))
}

func notUrlPath(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(" 404 path not found"))
}

func hi(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	w.Write([]byte("<h1> Hi " + name + "</h1>"))
}

func main() {

	handler := http.NewServeMux()
	handler.HandleFunc("/{$}", home) // Restrict this route to exact matches on / only.
	handler.HandleFunc("/", notUrlPath)
	handler.HandleFunc("/create", create)
	handler.HandleFunc("/view/{id}", view)
	handler.HandleFunc("/hi/{name}", hi)
	//don’t end in a trailing slash inless is a subtree path pattern

	log.Println(" server hosted on port 4000", handler)

	//err := http.ListenAndServe( ":4000", handler)
	err := http.ListenAndServe(":http-alt", handler) // 8080

	log.Fatal(err)

}
