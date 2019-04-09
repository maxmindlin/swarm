package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maxmindlin/swarm/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/articles", api.CreateArticleEndpoint).Methods("POST")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
