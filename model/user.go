package model

// User Table 전체
type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	LoginID    string `json:"loginId"`
	CreateDate string `json:"createDate"`
	ModifyDate string `json:"modifyDate"`
}

// UserInfo Front 유저에게 필요한 값
type UserInfo struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	LoginID string `json:"loginId"`
}

// Users 유저 배열 객체
var (
	Users       = map[int]*User{}
	UserSeq int = 1
)
