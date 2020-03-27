// 패키지명이 반드시 main인 파일만 실행파일로 생성 및 run 이 가능하다
package main

import (
	"strconv"

	"cndf.order.was/model"
	"cndf.order.was/route"
	"github.com/labstack/echo"
)

// 기본 포트 81
const PORT = 81

var messageId int64 = 0

// createMessage 더미 메시지 생성
func createMessage(message string, sender string) model.Message {
	messageId++
	return model.Message{
		ID:      messageId,
		Sender:  sender,
		Message: message,
	}
}

// main 프로그램의 시작점
func main() {

	e := echo.New()

	e.Logger.Info("Attempting to start HTTP Server.")

	e.GET("/", route.HelloWorld)
	e.GET("/user", route.UserList)
	e.GET("/user/:id", route.UserGet)
	e.PUT("/user", route.UserPut)
	e.DELETE("/user/:id", route.UserDelete)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))

	// // 특정 url주소로 Request가 들어올 경우 처리하는 handler함수를 할당 HandleRequest
	// http.HandleFunc("/message", httphandlers.HandleRequest)

	// // http WAS 서비스 온!
	// var err = http.ListenAndServe(":"+strconv.Itoa(PORT), nil)

	// // 로그 작성 후 main 함수는 끗 이후에는 http를 통해 request가 오면 callback 하는 보통의 WAS와 동일
	// if err != nil {
	// 	log.Panicln("Server failed starting. Error: ", err)
	// } else {
	// 	log.Panicln("Server starting. At ", PORT, " Port")
	// }
}
