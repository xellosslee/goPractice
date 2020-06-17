package storage

import (
	"database/sql"

	// using mysql
	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB 기본 DB 연결 시작
func ConnectDB() (*sql.DB, error) {
	// db, err := sql.Open("mysql", "blackk:cndfactory@)20@tcp(192.168.110.75:3306)/cndf.order?timeout=30s&charset=utf8mb4")
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/shoppingmall?timeout=30s&charset=utf8mb4")
	if err != nil {
		return nil, err
	}
	// 대기 커넥션 수를 10개로 설정
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	return db, nil
}

// SelectOne 하나의 행만 조회 할 때 사용
func SelectOne(db *sql.DB, query string, args ...interface{}) *sql.Row {
	conn := db.QueryRow(query, args...)
	return conn
}

// Select 조회 쿼리를 수행 할 때 사용
func Select(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	conn, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SelectPaging 페이지 처리를 하는 조회 쿼리를 수행 할 때 사용
// pagingType이 limit 이거나 공백이면 args 끝에 pageStart, pageCount 두개의 값은 필수
func SelectPaging(db *sql.DB, query string, pagingType string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if pagingType == "limit" || pagingType == "" {
		rows, err = db.Query("SELECT * FROM ("+query+") t LIMIT ?, ?", args...)
	} else {
		rows, err = db.Query("SELECT * FROM ("+query+") t LIMIT ?", args...)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Execute 실행류의 쿼리를 수행 할 때 사용(INSERT, UPDATE, DELETE)
func Execute(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
