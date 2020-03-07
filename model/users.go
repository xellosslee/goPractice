package model

// Users 유저정보 객체
type Users struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	LoginID    string `json:"loginId"`
	CreateDate string `json:"createDate"`
	ModifyDate string `json:"modifyDate"`
}
