package model

// 에러 메시지 정의
const (
	ErrMsgWorngParam   string = "잘못된 파라미터가 전달되었습니다"
	ErrMsgDbConnection string = "DB 연결에 실패하였습니다"
	ErrMsgCannotFound  string = "찾으시는 내용이 없습니다"
	ErrMsgProcFail     string = "@에 실패하였습니다"
	ErrMsgDuplicate    string = "동일한 @가 존재합니다"
	ErrMsgNotExists    string = "해당 @가 없습니다"
)
