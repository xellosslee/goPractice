package storage_test

import (
	"fmt"
	"log"
	"testing"

	"cndf.order.was/model"
	"cndf.order.was/storage"
)

/* 테스트에 사용된 테이블
DROP TABLE users;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `login_id` varchar(100) NOT NULL,
  `name` varchar(50) NOT NULL,
  `create_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='유저정보 테이블';
*/

func TestConnectDB(t *testing.T) {

	db := storage.ConnectDB()
	// Simple CRUD Test
	rows := storage.Query(db, "INSERT INTO users (name, login_id) VALUES(?,?)", "관리자", "gslee")
	// 적용된 rows 개수를 가져와서 0건이면 에러
	// 	t.Error("Wrong result")

	var id int
	var name, login_id string
	var users []model.Users
	rows = storage.Query(db, "SELECT id, name, login_id FROM users WHERE login_id = ?", "gslee")
	// 조회결과가 0건이면 에러
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
			// Update 결과가 0건이면 에러
		} else {
			// 나머지 row는 제거
			storage.Query(db, "DELETE FROM users WHERE id = ?", v.ID)
			// Delete 수행 했는데 결과가 0건이면 에러 - delete 할 내용이 없는 경우 에러로 체크 될 일 없음
		}
	}

	rows.Close()
}
