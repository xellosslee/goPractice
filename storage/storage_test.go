package storage_test

import (
	"fmt"
	"log"
	"testing"

	"cndf.order.was/model"
	"cndf.order.was/storage"
)

func TestConnectDB(t *testing.T) {

	db := storage.ConnectDB()
	// Simple CRUD Test
	rows := storage.Query(db, "INSERT INTO users (name, login_id) VALUES(?,?)", "관리자", "gslee")
	for rows.Next() {
		fmt.Println(rows.Scan())
	}

	var id int
	var name, login_id string
	var users []model.Users
	rows = storage.Query(db, "SELECT id, name, login_id FROM users WHERE login_id = ?", "gslee")
	for rows.Next() {
		err := rows.Scan(&id, &name, &login_id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, login_id)
		// pk값(id)을 뽑아서 배열에 저장
		users = append(users, model.Users{ID: id, Name: name, LoginID: login_id})
	}

	for i, v := range users {
		if i == 0 {
			// 첫번째 row는 name값을 변경
			storage.Query(db, "UPDATE users SET Name = ? WHERE id = ?", "test", v.ID)
		} else {
			// 나머지 row는 제거
			storage.Query(db, "DELETE FROM users WHERE id = ?", v.ID)
		}
	}

	rows.Close()
}
