package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gs.lee.was/sqlxDB"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	logging "github.com/op/go-logging"
)

type configData struct {
	serverName    string // 서버이름
	RunMode       string // release (서비스용) , debug (디버그용), local (로컬)
	Ssluse        string // Y,N
	Sslkey        string // ssl key 파일
	Sslcrt        string // ssl crt 파일
	MysqlConnInfo sqlxDB.DBConnectInfo
	HttpConfig    httpConfig
	CookieConfig  CookieConfig // cookie 관련 정보!1 - 수정하고 싶다면 수정하셈 !
	EmailInfo     emailConfig
	ButlerApiUrl  string   // 버틀러API 서버주소
	LogFileAccess *os.File // LogFile access
	LogFileInfo   *os.File // LogFile Info
	LogFileDebug  *os.File // LogFile 디버그
}

type emailConfig struct {
	Address  string
	Id       string
	Password string
	Host     string
	Port     int
}

var ServerConfig configData
var Log *logging.Logger
var Echo *echo.Echo

var MysqlObj *sqlxDB.SQLXforMysql //메인 db

func GetMysqlDB() *sqlx.DB {
	return MysqlObj.GetDBConn()
}

// 시스템 세팅 값 저장!!
func (ty *configData) InitConfig(cfPath string) error {

	er := ty.loadConfig(cfPath)

	if er != nil {
		return er
	}

	// mysql 등록
	MysqlObj = new(sqlxDB.SQLXforMysql) //메인 db
	MysqlObj.InitDB(ServerConfig.MysqlConnInfo, 10, 10, 300)
	_, err := MysqlObj.ConnectDB()
	if err != nil {
		return err
	}
	// http settting
	ty.InitHttp()

	// 로그 시작 로딩
	ty.loadLogInfo()

	return er
}

// config 파일 로드
func (ty *configData) loadConfig(cfPath string) error {

	b, err := ioutil.ReadFile(cfPath)
	if err != nil {
		log.Warn("Warn", "config file Not found", "mfconfig.json")
		return err
	}

	er := json.Unmarshal(b, &ServerConfig)
	if er != nil {
		log.Error("Error", "설정로드에러", er.Error())
		return er
	}

	return nil
}

// https 설정할때 쓰면 됨
func (ty *configData) InitHttp() {
	Echo = echo.New() //  echo 초기화
	if ServerConfig.RunMode == "debug" || ServerConfig.RunMode == "local" {
		Echo.Debug = true
	}
	Echo.Use(middleware.CORS()) // 추후 개발, 운영 url주소가 정해지면 해당 주소로만 접근 되도록 변경 필요
	// Echo.Use(session.Middleware(sessions.NewCookieStore([]byte("butler")))) // session store생성함
	// Echo.Use(middleware.Logger())  // fmt.Println 로그 (이거 호출하면 위의 loadLogInfo 에서 지정한 내용이 사라짐)
	Echo.Use(middleware.Recover()) // 에러 나면 다시 살려주는애인데. .못살릴 수도 있음.

	// // 기본 rest 통신 아이디/암호
	// Echo.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	// Be careful to use constant time comparison to prevent timing attacks
	// 	if subtle.ConstantTimeCompare([]byte(username), []byte("cndf")) == 1 &&
	// 		subtle.ConstantTimeCompare([]byte(password), []byte("actory")) == 1 {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))
	Echo.Validator = &CustomValidator{validator: validator.New()} // validator 체크 - 꼭 있어야되는애, tag랑 같이 쓰이는 듯
}

// 로그 설정
func (ty *configData) loadLogInfo() {
	var log = logging.MustGetLogger("gs.lee.was")
	var format = logging.MustStringFormatter(
		`%{color}[%{time:2006-01-02 15:04:05.000}] %{shortfile} - %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	Log = log

	var loggers []logging.Backend

	// 로그 폴더 없으면 생성
	os.Mkdir("./log", 0755)
	// 업로드 받을 폴더 생성
	os.MkdirAll("./upload/item", 0755)
	os.MkdirAll("./upload/product", 0755)
	os.MkdirAll("./upload/category", 0755)

	if ServerConfig.RunMode == "debug" || ServerConfig.RunMode == "local" {
		// log at debug file
		debugLog, err := os.OpenFile("log/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		debugLogBackend := logging.NewLogBackend(debugLog, "", 0)
		debugLogFormatter := logging.NewBackendFormatter(debugLogBackend, format)
		debugLogLevel := logging.AddModuleLevel(debugLogFormatter)
		debugLogLevel.SetLevel(logging.DEBUG, "")
		loggers = append(loggers, debugLogLevel)

		ServerConfig.LogFileDebug = debugLog
	}
	// log at info file
	infoLog, err := os.OpenFile("log/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	infoLogBackend := logging.NewLogBackend(infoLog, "", 0)
	infoFormatter := logging.NewBackendFormatter(infoLogBackend, format)
	infoLevel := logging.AddModuleLevel(infoFormatter)
	infoLevel.SetLevel(logging.INFO, "")
	loggers = append(loggers, infoLevel)
	ServerConfig.LogFileInfo = infoLog

	// backend1 에러 상황에서 화면에 표시되는 경우 출력
	standardOutput := logging.NewLogBackend(os.Stderr, "", 0)
	standardFormatter := logging.NewBackendFormatter(standardOutput, format)
	standardLog := logging.AddModuleLevel(standardFormatter)
	if ServerConfig.RunMode == "debug" || ServerConfig.RunMode == "local" {
		standardLog.SetLevel(logging.DEBUG, "")
	} else {
		standardLog.SetLevel(logging.ERROR, "")
	}
	loggers = append(loggers, standardLog)

	logging.SetBackend(loggers...)
	// echo 미들웨어에서 사용될 access 로그파일 오픈
	fpLog, err := os.OpenFile("log/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	ServerConfig.LogFileAccess = fpLog

	// access 에 적용될 내용들
	Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
		Output: fpLog,
	}))
}
