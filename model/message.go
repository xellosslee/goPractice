package model

// Message 커뮤니케이션용 객체
type Message struct {
	ID      int    `json:"id,omitempty"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// Users 유저 배열 객체
var (
	Messages   = map[int]*Message{}
	MessageSeq = 1
)
