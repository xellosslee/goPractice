package model

// Wait Table 전체
type Wait struct {
	ID         int64  `json:"id"`
	Create     string `json:"create"`
	ModifyDate string `json:"modify"`
}

// 현재 번호표 순번
type WaitPosition struct {
	ID int64 `json:"id"`
}
