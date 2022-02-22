package model

import (
	"gorm_page/database"

	"gorm.io/gorm"
)

type CityQueryInfo struct {
	PageInfo
	CountryCode string `json:"countryCode"`
	District    string `json:"district"`
}

type City struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
	District    string `json:"district"`
	Population  int    `json:"population"`
}

func (c *City) TableName() string {
	return "city"
}

func (c *City) SelectPageList(p *database.Page, wrapper map[string]interface{}) error {
	err := database.SelectPage(p, wrapper, City{})
	return err
}

func (c *City) CountAll(wrapper map[string]interface{}) int64 {
	var total int64
	database.DB.Model(&City{}).Where(wrapper).Count(&total)
	return total
}

func (c *City) SelectList(p *database.Page, wrapper map[string]interface{}) error {
	list := []City{}
	if err := database.DB.Model(&City{}).Scopes(Paginate(p)).Where(wrapper).Find(&list).Error; err != nil {
		return err
	}
	p.Data = list
	return nil
}

func Paginate(page *database.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.CurrentPage <= 0 {
			page.CurrentPage = 0
		} // 当前页小于0则置为0
		switch {
		case page.PageSize > 100:
			page.PageSize = 100
		case page.PageSize <= 0:
			page.PageSize = 10
		} // 限制size大小
		page.Pages = page.Total / page.PageSize
		if page.Total%page.PageSize != 0 {
			page.Pages++
		} // 计算总页数
		p := page.CurrentPage
		if page.CurrentPage > page.Pages {
			p = page.Pages
		} // 若当前页大于总页数则使用总页数
		size := page.PageSize
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size)) // 设置limit和offset
	}
}
