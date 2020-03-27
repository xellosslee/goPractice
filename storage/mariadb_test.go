package storage_test

import (
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
	res := storage.Execute(db, "INSERT INTO users (name, login_id) VALUES(?,?)", "관리자", "gslee")
	// 적용된 res 개수를 가져와서 0건이면 에러
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 0 {
		t.Error("Insert RowsAffected is 0")
	}
	insertId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	if insertId == 0 {
		t.Error("Insert pk is 0")
	}

	log.Println("insertId ", insertId)

	// 방금 추가한 레코드 한건 PK 값으로 조회
	var id int64
	var name, login_id string
	row := storage.SelectOne(db, "SELECT id, name, login_id FROM users WHERE id = ?", insertId)
	err = row.Scan(&id, &name, &login_id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SelectOne ", id, "_", name, "_", login_id)

	// users 테이블 전체 조회
	var users []model.User
	rows := storage.Select(db, "SELECT id, name, login_id FROM users")
	for rows.Next() {
		err := rows.Scan(&id, &name, &login_id)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Select Users ", id, "_", name, "_", login_id)
		// 배열에 저장
		users = append(users, model.User{ID: id, Name: name, LoginID: login_id})
	}

	// 방금 추가한 항목의 이름을 test로 변경
	res = storage.Execute(db, "UPDATE users SET Name = ? WHERE id = ?", "test", insertId)
	cnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 0 {
		t.Error("Update result count is 0")
	}
	// 업데이트 결과를 다시 조회해서 로그 적음
	row = storage.SelectOne(db, "SELECT id, name, login_id FROM users WHERE id = ?", insertId)
	err = row.Scan(&id, &name, &login_id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SelectOne ", id, "_", name, "_", login_id)

	// DELETE 수행
	res = storage.Execute(db, "DELETE FROM users WHERE id = ?", insertId)
	cnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 0 {
		t.Error("Delete result count is 0")
	}

	defer rows.Close()
}
