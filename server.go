package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/controller"
	"log"
	"net/http"
)

type Server struct {
	db     *sql.DB
	router *mux.Router
}

func New() *Server {
	return &Server{}
}

func (s *Server) Init() {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/ino")
	if err != nil {
		return
	}
	s.db = db
	s.router = s.Route()
}

func (s *Server) Route() *mux.Router {
	r := mux.NewRouter()

	user := &controller.User{s.db}
	_ = &controller.Reserve{s.db}

	// HelloWorld
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})
	r.HandleFunc("/reservation", controller.MockHandler).Methods("POST")

	// MySQL接続テスト
	r.HandleFunc("/user", middle(user.Get)).Methods("GET")
	r.HandleFunc("/user", middle(user.Post)).Methods("POST")

	return r
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":6000", s.router))
}
