package model

// Page 조회용
type Page struct {
	PageType string `json:"pageType"`
	Page     int    `json:"page"`
	Num      int    `json:"num"`
	ID       int    `json:"id"`
}

// PageResult Paging목록인 경우 해당 객체로 wrapping 해서 전달 필요
type PageResult struct {
	PageInfo PageInfo    `json:"pageInfo"`
	Result   interface{} `json:"result"`
}

// PageInfo 페이징에 필요한 추가 값 객체
type PageInfo struct {
	TotalPages  int `json:"totalPages"`
	TotalCounts int `json:"totalCounts"`
}
