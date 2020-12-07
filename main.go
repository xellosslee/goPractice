// 패키지명이 반드시 main인 파일만 실행파일로 생성 및 run 이 가능하다
package main
// 테스트
import (
	"net/http"
	"os"

	"cndf.order.was/route"
	"cndf.order.was/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cndf.order.was")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// main 프로그램의 시작점
func main() {

	e := echo.New()

	initLogConfig(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Cndfactory go web application server")
	})

	route.SetUserRouters(e)

	storage.ConnectDB()

	e.Logger.Fatal(e.Start(":80"))
}

func initLogConfig(e *echo.Echo) {
	debugLog, err := os.OpenFile("log/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer debugLog.Close()

	infoLog, err := os.OpenFile("log/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer infoLog.Close()

	// backend1 에러 상황에서 화면에 표시되는 경우 출력
	standardOutput := logging.NewLogBackend(os.Stderr, "", 0)
	infoLogBackend := logging.NewLogBackend(infoLog, "", 0)
	debugLogBackend := logging.NewLogBackend(debugLog, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	standardFormatter := logging.NewBackendFormatter(standardOutput, format)
	infoFormatter := logging.NewBackendFormatter(infoLogBackend, format)
	debugLogFormatter := logging.NewBackendFormatter(debugLogBackend, format)

	// Only errors and more severe messages should be sent to backend1
	standardLog := logging.AddModuleLevel(standardFormatter)
	standardLog.SetLevel(logging.ERROR, "")

	infoLevel := logging.AddModuleLevel(infoFormatter)
	infoLevel.SetLevel(logging.INFO, "")

	// Set the backends to be used.
	logging.SetBackend(standardLog, infoLevel, debugLogFormatter)

	// log.Debugf("debug")
	// log.Info("info")
	// log.Notice("notice")
	// log.Warning("warning")
	// log.Error("err")
	// log.Critical("crit")

	// echo 미들웨어에서 사용될 access 로그파일 오픈
	fpLog, err := os.OpenFile("log/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
		Output: fpLog,
	}))
}
