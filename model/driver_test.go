package model

import (
	"fmt"
	"github.com/prog470dev/inori-backend/db"
	"testing"
)

//TODO: modelの関数ごとにテスト関数を分割（その場合dbの使い回しはどうやるのか？）

func TestModelDriver(t *testing.T) {
	conf := &db.Config{}
	dbx, err := conf.Open("../dbconfig.yml") // testの場合は対象ファイルからの相対パス
	if err != nil {
		t.Fatalf("failed %s", err)
	}

	fmt.Println(dbx)

	_, err = dbx.Exec("show databases;")
	if err != nil {
		t.Fatalf("failed %s", err)
	}

	_, err = dbx.Exec("use ino;")
	if err != nil {
		t.Fatalf("failed %s", err)
	}

	//
	//driver := &Driver{
	//	ID:        0,
	//	FirstName: "Yamada",
	//	LastName:  "Taro",
	//	Grade:     "学部4年",
	//	Major:     "電子情報学類",
	//	Mail:      "foo@sample.com",
	//	Phone:     "00000000000",
	//	CarColor:  "赤",
	//	CarNumber: "0000",
	//}
	//
	//// 挿入
	//result, err := driver.Insert(dbx)
	//if err != nil {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
	//insertedID, err := result.LastInsertId()
	//if err != nil {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
	//
	//// 更新
	//driver.ID = insertedID
	//driver.CarColor = "青"
	//_, err = driver.Update(dbx)
	//if err != nil {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
	//
	//// 取得（更新確認）
	//givenDriver, err := DriverOne(dbx, insertedID)
	//if err != nil {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
	//if givenDriver.CarColor != "青" {
	//	t.Fatalf("expected: %s, actual: %s", "青", givenDriver.CarColor)
	//}
	//
	//// 取得（メアド）
	//givenDriver, err = DriverOneWithMail(dbx, "foo@sample.com")
	//if err != nil {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
	//
	//// 削除
	//result, err = driver.Delete(dbx)
	//if err != nil {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
	//givenDriver, err = DriverOne(dbx, insertedID)
	//if err != sql.ErrNoRows {
	//	t.Fatalf("failed %s", err)
	//	return
	//}
}
