package model

import "gorm_page_generic/database"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type PageResponse[T any] struct {
	CurrentPage int64 `json:"currentPage"`
	PageSize    int64 `json:"pageSize"`
	Total       int64 `json:"total"`
	Pages       int64 `json:"pages"` // 总页数
	Data        []T   `json:"data"`
}

func NewPageResponse[T any](page *database.Page[T]) *PageResponse[T] {
	return &PageResponse[T]{
		CurrentPage: page.CurrentPage,
		PageSize:    page.PageSize,
		Total:       page.Total,
		Pages:       page.Pages,
		Data:        page.Data,
	}
}
