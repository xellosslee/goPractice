package model

// Message 커뮤니케이션용 객체
type Message struct {
	ID      int    `json:"id,omitempty"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// Message 배열 및 pk 값 : 서버 메모리에 객체를 저장할 때만 사용됨 (테스트용)
var (
	Messages   = map[int]*Message{}
	MessageSeq = 1
)
