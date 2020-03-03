package main

import (
	"log"
	"net/http"
	"strconv"

	"cndf.order.was/httphandlers"
	"cndf.order.was/model"
	"cndf.order.was/storage"
)

const PORT = 8080

var messageId = 0

func createMessage(message string, sender string) model.Message {
	messageId++
	return model.Message{
		ID:      messageId,
		Sender:  sender,
		Message: message,
	}
}

func main() {
	log.Println("Creating dummy messages")

	storage.Add(createMessage("Testing", "1234"))
	storage.Add(createMessage("Testing Again", "5678"))
	storage.Add(createMessage("Testing A Third Time", "9012"))

	log.Println(storage.Get())

	log.Println("Attempting to start HTTP Server.")

	http.HandleFunc("/message", httphandlers.HandleRequest)

	var err = http.ListenAndServe(":"+strconv.Itoa(PORT), nil)

	if err != nil {
		log.Panicln("Server failed starting. Error: ", err)
	} else {
		log.Panicln("Server starting. At ", PORT, " Port")
	}
}
