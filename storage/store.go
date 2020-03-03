package storage

import (
	"cndf.order.was/model"
)

var store model.MessageList
var currentMaxID = 1

// Get 메모리에 담긴 메시지 목록 전달
func Get() model.MessageList {
	return store
}

// Add 메모리에 메시지 정보 추가
func Add(message model.Message) int {
	message.ID = currentMaxID
	currentMaxID++
	store = append(store, message)
	return message.ID
}

// Remove 메모리에 담긴 메시지 삭제
func Remove(id int) bool {
	index := -1

	for i, message := range store {
		if message.ID == id {
			index = i
		}
	}

	if index != -1 {
		store = append(store[:index], store[index+1:]...)
	}

	// Returns true if item was found & removed
	return index != -1
}
