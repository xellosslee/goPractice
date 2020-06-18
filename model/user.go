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

// User 배열 및 pk 값 : 서버 메모리에 객체를 저장할 때만 사용됨 (테스트용)
var (
	Users       = map[int]*User{}
	UserSeq int = 1
)
