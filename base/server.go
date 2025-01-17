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
	recommend := &controller.Recommend{s.db}

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

	// ライダーへのレコメンド通知
	go func() {
		//TODO: 定刻に実行されるように実装（今は定期的に時刻をチェックしている）
		t := time.NewTicker(60 * time.Minute) // 60毎にチェックして対象時刻の前後30分以内に入っているか確認
		for {
			<-t.C

			jst, _ := time.LoadLocation("Asia/Tokyo")
			schoolTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 0, 0, 0, jst)
			homeTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 15, 0, 0, 0, jst)

			// 登校方向（school）
			if time.Now().In(jst).After(schoolTime.Add(-30*time.Minute).Local()) &&
				time.Now().In(jst).Before(schoolTime.Add(30*time.Minute).Local()) {
				err := recommend.PushRecommend(0)
				if err != nil {
					log.Println(err)
				}
			}
			// 帰宅方向（home）
			if time.Now().In(jst).After(homeTime.Add(-30*time.Minute).Local()) &&
				time.Now().In(jst).Before(homeTime.Add(30*time.Minute).Local()) {
				err := recommend.PushRecommend(1)
				if err != nil {
					log.Println(err)
				}
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
	/** Driver: 画像設定 **/
	r.HandleFunc("/drivers/{driver_id:[0-9]+}/image", middle(driver.PostImage)).Methods("POST")

	/** Rider: サインアップ **/
	r.HandleFunc("/riders/signup", middle(rider.SignUpRider)).Methods("POST")
	/** Rider: サインイン **/
	r.HandleFunc("/riders/signin", middle(rider.SignInRider)).Methods("POST")
	/** Rider: 詳細 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}", middle(rider.GetRiderDetail)).Methods("GET")
	/** Rider: 更新 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}", middle(rider.UpdateRider)).Methods("PUT")
	/** Rider: 画像設定 **/
	r.HandleFunc("/riders/{rider_id:[0-9]+}/image", middle(rider.PostImage)).Methods("POST")

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

	/** Recommend: 強制実行 **/
	r.HandleFunc("/recommend/{dir:[a-z]+}", middle(recommend.ForcePushRecommend)).Methods("GET")

	return r
}

func (s *Server) Run() {
	inoPort := os.Getenv("INO_PORT")
	log.Fatal(http.ListenAndServe(":"+inoPort, s.router))

}
