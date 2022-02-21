package model

import (
	"gorm_page/database"
)

type PageInfo struct {
	CurrentPage int64 `json:"currentPage"`
	PageSize    int64 `json:"pageSize"`
}

type QueryInfo struct {
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

func (c *City) SelectPageList(p *database.Page, queryVo QueryInfo) error {
	p.CurrentPage = queryVo.CurrentPage
	p.PageSize = queryVo.PageSize
	wrapper := make(map[string]interface{}, 0)
	if queryVo.CountryCode != "" {
		wrapper["CountryCode"] = queryVo.CountryCode
	}
	if queryVo.District != "" {
		wrapper["District"] = queryVo.District
	}
	err := database.SelectPage(p, wrapper, City{})
	return err
}
