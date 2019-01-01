package base

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prog470dev/inori-backend/controller"
	"github.com/prog470dev/inori-backend/db"
	"github.com/prog470dev/inori-backend/model"
	"log"
	"net/http"
	"os"
	"time"
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
	token := &controller.Token{s.db}
	demand := &controller.Demand{s.db}

	// 需要集計処理の定期アップデート
	go func() {
		t := time.NewTicker(10 * time.Minute)
		for {
			<-t.C
			err := model.Aggregate(demand.DB)
			if err != nil {
				log.Println(err)
			}
		}
		t.Stop()
	}()

	// Health Check
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	/** Driver: サインアップ **/
	r.HandleFunc("/drivers/signup", middle(driver.SignUpDriver)).Methods("POST")
	/** Driver: サインイン **/
	r.HandleFunc("/drivers/signin", middle(driver.SignInDriver)).Methods("POST")
	/** Driver: 詳細 **/
	r.HandleFunc("/drivers/{driver_id:[0-9]+}", middle(driver.GetDriverDetail)).Methods("GET")
	/** Driver: 更新 **/
	r.HandleFunc("/drivers/{driver_id:[0-9]+}", middle(driver.UpdateDriver)).Methods("PUT")

	/** Rider: サインアップ **/
	r.HandleFunc("/riders/signup", middle(rider.SignUpRider)).Methods("POST")
	/** Rider: サインイン **/
	r.HandleFunc("/riders/signin", middle(rider.SignInRider)).Methods("POST")
	/** Rider: 詳細 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}", middle(rider.GetRiderDetail)).Methods("GET")
	/** Rider: 更新 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}", middle(rider.UpdateRider)).Methods("PUT")

	/** Offer: 追加**/
	r.HandleFunc("/offers", middle(offer.CreateOffer)).Methods("POST")
	/** Offer: 一覧（全て/選択ドライバ）**/
	r.HandleFunc("/offers", middle(offer.GetOffers)).Methods("GET") //クエリパラメータで渡す
	/** Offer: 詳細 **/
	r.HandleFunc("/offers/{offer_id:[0-9]+}", middle(offer.GetOfferDetail)).Methods("GET")
	/** Offer: 削除 **/
	r.HandleFunc("/offers/{offer_id:[0-9]+}", middle(offer.DeleteOffer)).Methods("DELETE")

	/**  Reservation: 一覧（選択ライダ）**/
	r.HandleFunc("/reservations", middle(reservation.GetRiderReservations)).Methods("GET") //予定変更 reservations を返す.
	/**  Reservation: 追加 **/
	r.HandleFunc("/reservations", middle(reservation.CreateReservation)).Methods("POST")
	/**  Reservation: 削除 **/
	r.HandleFunc("/reservations/{reservation_id:[0-9]+}", middle(reservation.CancelReservation)).Methods("DELETE")

	/** Token Push Driver **/
	r.HandleFunc("/tokens/push/drivers", middle(token.RegisterPushTokenDriver)).Methods("POST")
	/** Token Push Rider **/
	r.HandleFunc("/tokens/push/riders", middle(token.RegisterPushTokenRider)).Methods("POST")

	/** Demand: 重要集計取得 **/
	r.HandleFunc("/demand/aggregate/{dir:[a-z]+}", middle(demand.GetDemandAggregate)).Methods("GET")
	/** Demand: ライダーの需要取得 **/
	r.HandleFunc("/demand/{rider_id:[0-9]+}", middle(demand.GetDemandRider)).Methods("GET")
	/** Demand: ライダーの需要登録 **/
	r.HandleFunc("/demand", middle(demand.ResisterDemandRider)).Methods("POST")

	return r
}

func (s *Server) Run() {
	inoPort := os.Getenv("INO_PORT")
	log.Fatal(http.ListenAndServe(":"+inoPort, s.router))

}
