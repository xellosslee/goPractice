package storage_test

import (
	"database/sql"
	"log"
	"testing"

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

func TestPagination(t *testing.T) {

	db := storage.ConnectDB()
	var id int64
	var name, login_id string

	// users 테이블을 2개씩 각 페이지별로 조회
	// var pageNum int = 1
	// var getCount int = 10
	// var lastestId int64 = 0
	// var pagingType string = "num"
	var rows *sql.Rows
	// 최초 호출이어서 받았던 pk값이 없거나, 여러 페이지를 한꺼번에 이동할 때는 lastestId 를 0 으로 보내야 함 (ex: 1 => 8 page)
	var pageNum int = 1
	var getCount int = 10
	var lastestId int64 = 0
	var pagingType string = "id"
	// 페이징 처리가 1,2,3,4 등의 페이지가 아닌 보이지 않는 리스트를 전달하는 방식 (인스타, 페북 피드) 이라면 pageNum은 0 이고 lastestId & getCount 값을 보내야 함
	for {
		if pagingType == "id" {
			// id 가 0인 경우는 없으므로 최초 배열부터 가져옴
			// 반드시 정렬 순서가 정해져 있어야 함
			rows = storage.Select(db, "SELECT id, name, login_id FROM users WHERE id > ? ORDER BY id ASC LIMIT ?", lastestId, getCount)
		} else {
			rows = storage.Select(db, "SELECT id, name, login_id FROM users LIMIT ?, ?", (pageNum-1)*getCount, getCount)
		}
		var pageCnt = 0
		for rows.Next() {
			pageCnt++
			err := rows.Scan(&id, &name, &login_id)
			if err != nil {
				log.Fatal(err)
			}
			if pagingType == "id" {
				lastestId = id
				log.Println("Select Users Pagination By PrimaryKey Value ", pageNum, " user : ", id, "_", name, "_", login_id)
			} else {
				log.Println("Select Users Pagination By PageNum ", pageNum, " user : ", id, "_", name, "_", login_id)
			}
		}
		// 요청한 개수와 리턴된 개수가 다르다면 그냥 break 로 루프 종료
		if pageCnt != getCount {
			break
		} else {
			pageNum++
		}
	}

	defer rows.Close()
}
