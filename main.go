package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	log.Println("Server listening in http://localhost:4000/")
	err := http.ListenAndServe(":4000", mux)
	log.Println(err)

}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("showSnippet"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create Snippet"))
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("hello i want to sleep"))
}
