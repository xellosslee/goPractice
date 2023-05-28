package sqlxDB

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBConnectInfo struct {
	StrID     string
	StrPasswd string
	StrIP     string
	NPort     int
	StrDBname string
}

type MetalScanner struct {
	valid bool
	value interface{}
}

func DBQuerySelect(dbx *sqlx.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {

	ctx := context.Background()
	err := dbx.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	if dbx == nil {
		return nil, fmt.Errorf("Database not connection")
	}

	rows, err := dbx.Query(query, args...)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	var send []map[string]interface{}

	// iterate over each row
	for rows.Next() {

		tmp := make(map[string]interface{})
		row := make([]interface{}, len(columns))
		for idx := range columns {
			row[idx] = new(MetalScanner)
		}

		err := rows.Scan(row...)
		if err != nil {
			if rows != nil {
				rows.Close()
			}
			return nil, err
		}

		// 컬럼명 = 키 버전
		for idx, column := range columns {
			var scanner = row[idx].(*MetalScanner)
			strtmp := scanner.value
			tmp[column] = strtmp
		}
		// c1~c99 = 키 버전
		// for idx, _ := range columns {
		// 	var scanner = row[idx].(*MetalScanner)
		// 	strtmp := scanner.value
		// 	tmp["c"+strconv.Itoa(idx+1)] = strtmp
		// }

		send = append(send, tmp)
	}

	// check the error from rows
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return send, nil
}

func (scanner *MetalScanner) getBytes(src interface{}) []byte {
	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}

func (scanner *MetalScanner) Scan(src interface{}) error {
	// t := src.(type)
	switch src.(type) {
	case int:
		if value, ok := src.(int); ok {
			scanner.value = value
			scanner.valid = true
		}
	case int64:
		if value, ok := src.(int64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case float64:
		if value, ok := src.(float64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case float32:
		if value, ok := src.(float32); ok {
			scanner.value = value
			scanner.valid = true
		}
	case bool:
		if value, ok := src.(bool); ok {
			scanner.value = value
			scanner.valid = true
		}
	case string:
		scanner.value = string(src.(string))
		scanner.valid = true
	case []byte:
		value := scanner.getBytes(src)
		scanner.value = string(value)
		scanner.valid = true
	case time.Time:
		if value, ok := src.(time.Time); ok {
			scanner.value = value
			scanner.valid = true
		}
	case nil:
		scanner.value = nil
		scanner.valid = true
	}
	return nil
}
