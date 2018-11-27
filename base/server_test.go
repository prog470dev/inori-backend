package base

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func defaultServer() *Server {
	s := New()
	s.Init("../db_config.yml") // testの場合は対象ファイルからの相対パス
	return s
}

func TestHealthCheck(t *testing.T) {
	s := defaultServer()

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatalf("ping failed %s", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if string(b) != "pong" {
		t.Fatalf("expected: %s, actual: %s", "pong", b)
	}
	defer resp.Body.Close()
}

// テスト前にかならずDBを初期化(IDの関係)

//func TestSystemDriver(t *testing.T) {
//	s := defaultServer()
//
//	ts := httptest.NewServer(s.router)
//	defer ts.Close()
//
//	client := new(http.Client)
//
//	/*
//		サンプルドライバ
//
//			ID:        1,
//			FirstName: "Yamada",
//			LastName:  "Taro",
//			Grade:     "学部4年",
//			Major:     "電子情報学類",
//			Mail:      "foo@sample.com",
//			Phone:     "00000000000",
//			CarColor:  "赤",
//			CarNumber: "0000",
//	*/
//
//	cases := []struct {
//		method          string
//		path            string
//		body            string
//		expected_status int
//		expected_body   string
//	}{
//		// サインアップ
//		{"POST", "/drivers/singup", `{"id":0,"first_name":"Yamada","last_name":"Taro","grade":"学部4年","major":"電子情報学類","mail":"foo@sample.com","phone":"00000000000","car_color":"赤","car_number":"0000"}`, 200, `{"id":1}`},
//		// サインイン
//		{"POST", "/drivers/singin", `{"mail":"foo@sample.com"}`, 200, `{"driver":{"id":1,"first_name":"Yamada","last_name":"Taro","grade":"学部4年","major":"電子情報学類","mail":"foo@sample.com","phone":"00000000000","car_color":"赤","car_number":"0000"}}`},
//		// 詳細取得
//		{"GET", "/drivers/1", ``, 200, `{"driver":{"id":1,"first_name":"Yamada","last_name":"Taro","grade":"学部4年","major":"電子情報学類","mail":"foo@sample.com","phone":"00000000000","car_color":"赤","car_number":"0000"}}`},
//		// 更新
//		{"PUT", "/drivers/1", `{"id":1,"first_name":"Yamada","last_name":"Taro","grade":"学部4年","major":"電子情報学類","mail":"foo@sample.com","phone":"00000000000","car_color":"赤","car_number":"0000"}`, 200, `{"id":1}`},
//	}
//
//	for i, c := range cases {
//		method := c.method
//		path := ts.URL + c.path
//		body := c.body
//		expected_status := c.expected_status
//		expected_body := c.expected_body
//
//		req, err := http.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
//		if err != nil {
//			t.Fatalf("err %d\n", i)
//			return
//		}
//
//		req.Header.Set("Content-Type", "application/json")
//
//		resp, err := client.Do(req)
//		if err != nil {
//			t.Fatalf("err %d\n", i)
//			return
//		}
//		defer resp.Body.Close()
//
//		//ステータスコード確認
//		if resp.StatusCode != expected_status {
//			t.Fatalf("expected: %d, actual: %d (case: %d)", expected_status, resp.StatusCode, i)
//			return
//		}
//
//		//レスポンスBODY取得
//		rbody, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			t.Error(resp.StatusCode)
//			return
//		}
//
//		actual_body := strings.Trim(string(rbody), "\n")
//
//		//レスポンスBODY確認
//		if actual_body != expected_body {
//			t.Fatalf("expected: %s, actual: %s (case: %d)", expected_body, actual_body, i)
//			return
//		}
//	}
//}
