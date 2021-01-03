package mysql

import (
	"database/sql"
	"log"

	// using mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DB Connection Pool 적용된 글로벌 변수
var DB *sql.DB

// ConnectDB 기본 DB 연결 시작
func ConnectDB() {
	var err error
	DB, err = sql.Open("mysql", "root@tcp(localhost:3306)/shoppingmall?timeout=30s&charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	// 대기 커넥션 수를 10개로 설정
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(10)
}

// SelectOne 하나의 행만 조회 할 때 사용
func SelectOne(query string, args ...interface{}) *sql.Row {
	conn := DB.QueryRow(query, args...)
	return conn
}

// Select 조회 쿼리를 수행 할 때 사용
func Select(query string, args ...interface{}) (*sql.Rows, error) {
	conn, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SelectPaging 페이지 처리를 하는 조회 쿼리를 수행 할 때 사용
// pagingType이 limit 이거나 공백이면 args 끝에 pageStart, pageCount 두개의 값은 필수
func SelectPaging(query string, pagingType string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if pagingType == "limit" || pagingType == "" {
		rows, err = DB.Query("SELECT * FROM ("+query+") t LIMIT ?, ?", args...)
	} else {
		rows, err = DB.Query("SELECT * FROM ("+query+") t LIMIT ?", args...)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Execute 실행류의 쿼리를 수행 할 때 사용(INSERT, UPDATE, DELETE)
func Execute(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
