package storage

import (
	"cndf.order.was/model"
)

// 이 패키지는 메모리에 저장할 데이터가 필요한 경우 활용 될 수 있음

var store model.MessageList
var currentMaxID int64 = 1

// Get 메모리에 담긴 메시지 목록 전달
func Get() model.MessageList {
	return store
}

// Add 메모리에 메시지 정보 추가
func Add(message model.Message) int64 {
	message.ID = currentMaxID
	currentMaxID++
	store = append(store, message)
	return message.ID
}

// Remove 메모리에 담긴 메시지 삭제
func Remove(id int64) bool {
	index := -1

	// foreach 구문과 비슷하게 loop 수행
	for i, message := range store {
		if message.ID == id {
			index = i
		}
	}

	if index != -1 {
		// 배열에 추가 배열을 만들어서 합치는 append함수와
		// 배열의 일부만 가져오는 slice 문법을 활용하여
		// store 리스트객체 내에서 삭제될 index 앞뒤의 Array를 새로운 Array로 조합
		store = append(store[:index], store[index+1:]...)
	}

	// Returns true if item was found & removed
	return index != -1
}
