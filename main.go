package main

import (
	"log"
	"net/http"
)

func Handler_foo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("foo"))
}

func Handler_bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("bar"))
}

// type Server struct {
// 	Mux *http.ServeMux
// }

// var Srv Server = Server{}

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/h1", Handler_foo)
	mux.HandleFunc("/h2", Handler_bar)
	return mux
}

func main() {
	mux := setupRouter()
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
