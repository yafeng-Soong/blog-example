package database

type Page struct {
	CurrentPage int64
	PageSize    int64
	Total       int64
	Pages       int64
	Data        interface{}
}
