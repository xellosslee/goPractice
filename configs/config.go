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

var debugmodeflg bool //디버그 모드 플레그 디버그:true,

// 서버 버전 정보 - in config.json
type VersionInfo struct {
	CUR_VER  string
	PACE_VER string
	JC1_VER  string
}

type DXConfig struct {
	TRANS_LOG_YN string // 작업로그 쓸지 말지 결정하기 !!  DB에 저장을 하느냐 마느냐 .
}

type EmailConfig struct {
	Address  string
	Id       string
	Password string
	Host     string
	Port     int
}

type ConfigData struct {
	ServerName    string // 서버이름
	RunMode       string // release (서비스용) , debug (디버그용), local (로컬)
	Ssluse        string // Y,N
	Sslkey        string // ssl key 파일
	Sslcrt        string // ssl crt 파일
	MysqlConnInfo sqlxDB.DBConnectInfo
	HttpConfig    httpConfig
	CookieConfig  CookieConfig // cookie 관련 정보!1 - 수정하고 싶다면 수정하셈 !
	VersionInfo   VersionInfo
	DXConfig      DXConfig
	EmailInfo     EmailConfig
	ButlerApiUrl  string // 버틀러API 서버주소
}

type SingleConfig struct {
	cnf    ConfigData
	CnfLog StaticLog
	E      *echo.Echo
}

var ConfigPTR *SingleConfig
var MysqlObj *sqlxDB.SQLXforMysql //메인 db

func IsDebugmode() bool {
	return debugmodeflg
}

func SetDebugmode(b bool) {
	debugmodeflg = b
}

func GetConfig() *SingleConfig {

	if ConfigPTR == nil {
		ConfigPTR = new(SingleConfig)
	}
	return ConfigPTR
}

func GetConfigData() *ConfigData {
	return &ConfigPTR.cnf
}

func GetConfigLog() *logging.Logger {
	return ConfigPTR.CnfLog.Log
}

func GetMysqlDB() *sqlx.DB {
	return MysqlObj.GetDBConn()
}

func GetDXConfig() DXConfig {
	return ConfigPTR.cnf.DXConfig
}

// 시스템 세팅 값 저장!!
func (ty *SingleConfig) InitConfig(cfPath string) error {

	er := ty.loadConfig(cfPath)

	if er != nil {
		return er
	}

	// mysql 등록
	MysqlObj = new(sqlxDB.SQLXforMysql) //메인 db
	MysqlObj.InitDB(ty.cnf.MysqlConnInfo, 10, 10, 300)
	_, err := MysqlObj.ConnectDB()
	if err != nil {
		return err
	}

	// http settting
	ty.E = echo.New() //  echo 초기화

	if ty.cnf.RunMode == "debug" || ty.cnf.RunMode == "local" {
		ty.E.Debug = true
	}

	// 로그 시작 로딩
	ty.loadLogInfo()
	ty.InitHttp()

	return er
}

// config 파일 로드
func (ty *SingleConfig) loadConfig(cfPath string) error {

	b, err := ioutil.ReadFile(cfPath)
	if err != nil {
		log.Warn("Warn", "config file Not found", "mfconfig.json")
		return err
	}

	er := json.Unmarshal(b, &ty.cnf)
	if er != nil {
		log.Error("Error", "설정로드에러", er.Error())
		return er
	}

	if ty.cnf.RunMode == "debug" || ty.cnf.RunMode == "local" {
		SetDebugmode(true)
	} else {
		SetDebugmode(false)
	}

	return nil
}

// 로그 설정 로드
func (ty *SingleConfig) loadLogInfo() {
	var log = logging.MustGetLogger("gs.lee.was")
	var format = logging.MustStringFormatter(
		`%{color}[%{time:2006-01-02 15:04:05.000}] %{shortfile} - %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	ty.CnfLog.Log = log

	var loggers []logging.Backend

	// 로그 폴더 없으면 생성
	os.Mkdir("./log", 0755)
	// 업로드 받을 폴더 생성
	os.MkdirAll("./upload/item", 0755)
	os.MkdirAll("./upload/product", 0755)
	os.MkdirAll("./upload/category", 0755)

	if IsDebugmode() {
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

		ty.CnfLog.Defung = debugLog
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
	ty.CnfLog.InfoLog = infoLog

	// backend1 에러 상황에서 화면에 표시되는 경우 출력
	standardOutput := logging.NewLogBackend(os.Stderr, "", 0)
	standardFormatter := logging.NewBackendFormatter(standardOutput, format)
	standardLog := logging.AddModuleLevel(standardFormatter)
	if ty.cnf.RunMode == "debug" || ty.cnf.RunMode == "local" {
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

	ty.CnfLog.FpLog = fpLog // main 쪽에 가서 정리하셈

	// access 에 적용될 내용들
	ty.E.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
		Output: ty.CnfLog.FpLog,
	}))
}

// https 설정할때 쓰면 됨
func (ty *SingleConfig) InitHttp() {

	ty.E.Use(middleware.CORS()) // 추후 개발, 운영 url주소가 정해지면 해당 주소로만 접근 되도록 변경 필요
	// ty.E.Use(session.Middleware(sessions.NewCookieStore([]byte("butler")))) // session store생성함
	// ty.E.Use(middleware.Logger())  // fmt.Println 로그 (이거 호출하면 위의 loadLogInfo 에서 지정한 내용이 사라짐)
	ty.E.Use(middleware.Recover()) // 에러 나면 다시 살려주는애인데. .못살릴 수도 있음.

	// // 기본 rest 통신 아이디/암호
	// ty.E.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	// Be careful to use constant time comparison to prevent timing attacks
	// 	if subtle.ConstantTimeCompare([]byte(username), []byte("cndf")) == 1 &&
	// 		subtle.ConstantTimeCompare([]byte(password), []byte("actory")) == 1 {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	ty.E.Validator = &CustomValidator{validator: validator.New()} // validator 체크 - 꼭 있어야되는애, tag랑 같이 쓰이는 듯
}
