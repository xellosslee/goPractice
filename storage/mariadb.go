package storage

import (
	"database/sql"
	"log"

	// using mysql
	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB 기본 DB 연결 시작
func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "blackk:cndfactory@)20@tcp(192.168.110.75:3306)/cndf.order?timeout=30s&charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}
	// 대기 커넥션 수를 10개로 설정
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	return db
}

// Query 모든 쿼리를 직접 수행 할 때 사용
func Query(db *sql.DB, query string, args ...interface{}) *sql.Rows {
	conn, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
