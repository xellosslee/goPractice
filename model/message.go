package model

// Message 커뮤니케이션용 객체
type Message struct {
	ID      int    `json:"id,omitempty"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
