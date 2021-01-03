package oracle_test

import (
	"fmt"
	"log"
	"testing"

	"cndf.order.was/storage/oracle"
)

func TestOracleDB(t *testing.T) {

	oracle.ConnectDB()

	// Simple CRUD Test
	rows, err := oracle.Select("SELECT SERVICE_CD, EST_CD FROM AR_CASH_IN WHERE ROWNUM <= :1", 2)
	if err != nil {
		log.Fatal(err)
	}
	var serviceCd, estCd string
	for rows.Next() {
		rows.Scan(&serviceCd, &estCd)
		fmt.Println(serviceCd, estCd)
	}
	defer rows.Close()
}
