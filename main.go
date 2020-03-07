// 패키지명이 반드시 main인 파일만 실행파일로 생성 및 run 이 가능하다
package main

import (
	"log"
	"net/http"
	"strconv"

	"cndf.order.was/httphandlers"
	"cndf.order.was/model"
	"cndf.order.was/storage"
)

// 기본 포트는 8080포트
const PORT = 8080

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
	log.Println("Creating dummy messages")

	storage.Add(createMessage("Testing", "1234"))
	storage.Add(createMessage("Testing Again", "5678"))
	storage.Add(createMessage("Testing A Third Time", "9012"))

	log.Println(storage.Get())

	log.Println("Attempting to start HTTP Server.")

	// 특정 url주소로 Request가 들어올 경우 처리하는 handler함수를 할당 HandleRequest
	http.HandleFunc("/message", httphandlers.HandleRequest)

	// http WAS 서비스 온!
	var err = http.ListenAndServe(":"+strconv.Itoa(PORT), nil)

	// 로그 작성 후 main 함수는 끗 이후에는 http를 통해 request가 오면 callback 하는 보통의 WAS와 동일
	if err != nil {
		log.Panicln("Server failed starting. Error: ", err)
	} else {
		log.Panicln("Server starting. At ", PORT, " Port")
	}
}
