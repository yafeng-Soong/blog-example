package database

type Page struct {
	CurrentPage int64       `json:"currentPage"`
	PageSize    int64       `json:"pageSize"`
	Total       int64       `json:"total"`
	Pages       int64       `json:"pages"`
	Data        interface{} `json:"data"`
}
