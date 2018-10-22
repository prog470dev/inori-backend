package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/prog470dev/inori-backend/controller"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/reservation", controller.MockHandler).Methods("GET")

	http.Handle(":8080", r)
}
