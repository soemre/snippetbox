package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Handlers
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Serve
	log.Println("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatalln(err)
}
