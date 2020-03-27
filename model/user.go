package model

// User 유저정보 객체
type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	LoginID    string `json:"loginId"`
	CreateDate string `json:"createDate"`
	ModifyDate string `json:"modifyDate"`
}

// Users 유저 배열 객체
var (
	Users         = map[int64]*User{}
	UserSeq int64 = 1
)
