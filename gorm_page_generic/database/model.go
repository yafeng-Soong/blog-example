package database

import (
	"gorm.io/gorm"
)

type Page[T any] struct {
	CurrentPage int64
	PageSize    int64
	Total       int64
	Pages       int64
	Data        []T
}

func (page *Page[T]) SelectPage(wrapper map[string]interface{}) (e error) {
	e = nil
	var model T
	DB.Model(&model).Where(wrapper).Count(&page.Total)
	if page.Total == 0 {
		page.Data = []T{}
		return
	}
	e = DB.Model(&model).Where(wrapper).Scopes(Paginate(page)).Find(&page.Data).Error
	if e != nil {
		return
	}
	return
}

func Paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.CurrentPage <= 0 {
			page.CurrentPage = 0
		}
		switch {
		case page.PageSize > 100:
			page.PageSize = 100
		case page.PageSize <= 0:
			page.PageSize = 10
		}
		page.Pages = page.Total / page.PageSize
		if page.Total%page.PageSize != 0 {
			page.Pages++
		}
		p := page.CurrentPage
		if page.CurrentPage > page.Pages {
			p = page.Pages
		}
		size := page.PageSize
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}
