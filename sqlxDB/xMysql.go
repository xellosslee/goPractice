package sqlxDB

import (
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type SQLXforMysql struct {
	bInited         bool
	stDBConnectInfo DBConnectInfo //연결정보
	mDBconn         *sqlx.DB
	mMaxOpenConn    int
	mMaxIdelConn    int

	mJobTime time.Time
	idleTime int64
	closeFlg bool //데이터베이스 close됐는가
}

func (ty *SQLXforMysql) InitDB(tyDbInfoStruct DBConnectInfo, nMaxOpenConn int, nMaxIdelConn int, nIdleChkTime int64) {
	ty.stDBConnectInfo = tyDbInfoStruct
	ty.mMaxOpenConn = nMaxOpenConn
	ty.mMaxIdelConn = nMaxIdelConn

	ty.idleTime = nIdleChkTime
	ty.closeFlg = false
	if ty.idleTime < 60 {
		ty.idleTime = 60
	}
}

func (ty *SQLXforMysql) ConnectDB() (dbconn *sqlx.DB, err error) {

	// 접속경로 지정하는 부분 ex => "root:fnakfn100djr@(localhost:3306)/sakila"
	// port_list := strconv.Itoa(ty.stDBConnectInfo.NPort)
	// stringArray := []string{ty.stDBConnectInfo.StrID, ":", ty.stDBConnectInfo.StrPasswd, "@tcp(", ty.stDBConnectInfo.StrIP, ":", port_list, ")/", ty.stDBConnectInfo.StrDBname}
	// justString := strings.Join(stringArray, "")
	var connDB = ty.stDBConnectInfo.StrID + ":" + ty.stDBConnectInfo.StrPasswd +
		"@tcp(" + ty.stDBConnectInfo.StrIP + ":" + strconv.Itoa(ty.stDBConnectInfo.NPort) + ")/" + ty.stDBConnectInfo.StrDBname

	db, err := sqlx.Connect("mysql", connDB)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(ty.mMaxIdelConn)
	db.SetMaxOpenConns(ty.mMaxOpenConn)

	ty.mDBconn = db
	ty.bInited = true

	ty.mJobTime = time.Now()
	ty.closeFlg = false
	// go ty.dbIdleUnconnectionChecker() // DB connection 문제 인거 같은데. 이제 해결된듯,,
	return db, nil
}

func (ty *SQLXforMysql) GetDBConn() *sqlx.DB {
	return ty.mDBconn
}

/*
//dbIdleUnconnectionChecker 특정시간 이상 디비작업 없을경우 mysql이 커넥션을 끊는다. 이를 해결하기위해 마지막 작업타임을 저장후 특정시간 이상 작업이 없다면 의미없는 작업을 진행하여 연결유지하도록 한다.
func (ty *SQLXforMysql) dbIdleUnconnectionChecker() {
	for {
		if ty.closeFlg == true {
			return
		}

		time.Sleep(time.Duration(int64(time.Second)))

		now := time.Now()

		idT := ty.mJobTime.Add(time.Duration(int64(time.Second) * ty.idleTime))

		if idT.Unix() < now.Unix() {
			sql := "select"
			ty.DBQuerySelect(sql)
		}
	}
}
*/
