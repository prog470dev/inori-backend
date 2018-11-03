package base

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/controller"
	"github.com/prog470dev/inori-backend/db"
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

func (s *Server) Init(filename string) {
	conf := &db.Config{}
	dbx, err := conf.Open(filename)
	if err != nil {
		return
	}
	s.db = dbx
	s.router = s.Route()
}

func (s *Server) Route() *mux.Router {
	r := mux.NewRouter()

	driver := &controller.Driver{s.db}
	rider := &controller.Rider{s.db}
	offer := &controller.Offer{s.db}
	reservation := &controller.Reservation{s.db}

	// HelloWorld
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	/** Driver **/
	r.HandleFunc("/drivers/singup", middle(driver.SignUpDriver)).Methods("POST")
	r.HandleFunc("/drivers/singin", middle(driver.SignInDriver)).Methods("POST")
	r.HandleFunc("/drivers/{driver_id:[0-9]+}", middle(driver.GetDriverDetail)).Methods("GET")
	r.HandleFunc("/drivers/{driver_id:[0-9]+}", middle(driver.UpdateDriver)).Methods("PUT")

	// Rider
	r.HandleFunc("/riders/singup", middle(rider.SignUpRider)).Methods("POST")

	// Offer
	r.HandleFunc("/offers", middle(offer.GetDriverOffers)).Methods("GET") //クエリパラメータで渡す
	r.HandleFunc("/offers/{offer_id:[0-9]+}", middle(offer.GetOfferDetail)).Methods("GET")
	r.HandleFunc("/offers/{offer_id:[0-9]+}", middle(offer.DeleteOffer)).Methods("DELETE")
	r.HandleFunc("/offers", middle(offer.CreateOffer)).Methods("POST")

	// Reservation
	r.HandleFunc("/reservations", middle(reservation.CreateReservation)).Methods("POST")

	return r
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", s.router))
}
