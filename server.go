package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/prog470dev/inori-backend/controller"
	"fmt"
	"log"
)


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})
	r.HandleFunc("/reservation", controller.MockHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":6000", r))
}
