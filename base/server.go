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

	/** Driver: サインアップ **/
	r.HandleFunc("/drivers/singup", middle(driver.SignUpDriver)).Methods("POST")
	/** Driver: サインイン **/
	r.HandleFunc("/drivers/singin", middle(driver.SignInDriver)).Methods("POST")
	/** Driver: 詳細 **/
	r.HandleFunc("/drivers/{driver_id:[0-9]+}", middle(driver.GetDriverDetail)).Methods("GET")
	/** Driver: 更新 **/
	r.HandleFunc("/drivers/{driver_id:[0-9]+}", middle(driver.UpdateDriver)).Methods("PUT")

	/** Rider: サインアップ **/
	r.HandleFunc("/riders/singup", middle(rider.SignUpRider)).Methods("POST")
	/** Rider: サインイン **/
	r.HandleFunc("/riders/singin", middle(rider.SignInRider)).Methods("POST")
	/** Rider: 詳細 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}", middle(rider.GetRiderDetail)).Methods("GET")
	/** Rider: 更新 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}", middle(rider.UpdateRider)).Methods("PUT")

	/** Offer: 追加**/
	r.HandleFunc("/offers", middle(offer.CreateOffer)).Methods("POST")
	/** Offer: 一覧（全て/選択ドライバ）**/
	r.HandleFunc("/offers", middle(offer.GetDriverOffers)).Methods("GET") //クエリパラメータで渡す
	/** Offer: 詳細 **/
	r.HandleFunc("/offers/{offer_id:[0-9]+}", middle(offer.GetOfferDetail)).Methods("GET")
	/** Offer: 削除 **/
	r.HandleFunc("/offers/{offer_id:[0-9]+}", middle(offer.DeleteOffer)).Methods("DELETE")

	/**  Reservation: 一覧（選択ライダ）**/
	r.HandleFunc("/reservations", middle(reservation.GetRiderOffers)).Methods("GET") //予定変更 reservations を返す.
	/**  Reservation: 追加 **/
	r.HandleFunc("/reservations", middle(reservation.CreateReservation)).Methods("POST")
	/**  Reservation: 削除 **/
	r.HandleFunc("/reservations/{reservation_id:[0-9]+}", middle(reservation.CancelReservation)).Methods("DELETE")

	return r
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", s.router))
}