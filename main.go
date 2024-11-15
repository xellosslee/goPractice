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

	er := configs.ServerConfig.InitConfig(".env.json")
	if er != nil {
		fmt.Println("설정정보로드에러")
		os.Exit(1)
	}

	if configs.ServerConfig.RunMode != "local" {
		setSchedule()
	}

	route.SetUserRouters(configs.Echo)
	route.SetWaitRoute(configs.Echo)

	runtime.GOMAXPROCS(runtime.NumCPU())

	configData := configs.ServerConfig

	configs.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	configs.Echo.Static("/upload", "upload")

	if configData.Ssluse == "Y" {
		// 보안 접속을 위한 https 서버 실행
		// "cert.pem"과 "privkey.pem" 파일이 필요함
		configs.Echo.Logger.Fatal(configs.Echo.StartTLS(":"+strconv.Itoa(configData.HttpConfig.Port), configData.Sslcrt, configData.Sslkey))
	} else {
		// 일반적인 http 서버 실행
		configs.Echo.Logger.Fatal(configs.Echo.Start(":" + strconv.Itoa(configData.HttpConfig.Port)))
	}

	defer configs.ServerConfig.LogFileAccess.Close()
	defer configs.ServerConfig.LogFileInfo.Close()
	defer configs.ServerConfig.LogFileDebug.Close()
}

func setSchedule() {
	c := cron.New()
	// // 새벽 4시마다 수행
	// c.AddFunc("0 0 4 * * *", func() {
	// 	route.CheckStoreMagamStatus()
	// })
	c.Start()
}
