/*
http 통신 방법중 echo를 테스트 하기 위한것으로
사용하기 편하게 쪼개서 정리한다.
1. 라우터 정리
2. 실행 파일 정리
3. 테스트 DB 출력 및 등록
*/
package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/robfig/cron"

	echoSwagger "github.com/swaggo/echo-swagger"
	"gs.lee.was/configs"
	_ "gs.lee.was/docs"
	"gs.lee.was/route"
)

func main() {

	cnf := configs.GetConfig()
	er := cnf.InitConfig(".env.json")
	if er != nil {
		fmt.Println("설정정보로드에러")
		os.Exit(1)
	}

	if configs.GetConfigData().RunMode != "local" {
		setSchedule()
	}

	route.SetUserRouters(cnf.E)

	runtime.GOMAXPROCS(runtime.NumCPU())

	configData := configs.GetConfigData()

	cnf.E.GET("/swagger/*", echoSwagger.WrapHandler)
	cnf.E.Static("/upload", "upload")

	if configData.Ssluse == "Y" {
		// 보안 접속을 위한 https 서버 실행
		// "cert.pem"과 "privkey.pem" 파일이 필요함
		cnf.E.Logger.Fatal(cnf.E.StartTLS(":"+strconv.Itoa(configData.HttpConfig.Port), configData.Sslcrt, configData.Sslkey))
	} else {
		// 일반적인 http 서버 실행
		cnf.E.Logger.Fatal(cnf.E.Start(":" + strconv.Itoa(configData.HttpConfig.Port)))
	}

	defer cnf.CnfLog.Defung.Close()
	defer cnf.CnfLog.InfoLog.Close()
	defer cnf.CnfLog.FpLog.Close()
}

func setSchedule() {
	c := cron.New()
	// // 새벽 4시마다 수행
	// c.AddFunc("0 0 4 * * *", func() {
	// 	route.CheckStoreMagamStatus()
	// })
	c.Start()
}
