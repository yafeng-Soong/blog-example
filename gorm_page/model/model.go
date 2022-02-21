package model

import (
	"gorm_page/database"
)

type PageInfo struct {
	CurrentPage int64 `json:"currentPage"`
	PageSize    int64 `json:"pageSize"`
}

type CityQueryInfo struct {
	PageInfo
	CountryCode string `json:"countryCode"`
	District    string `json:"district"`
}

type CountryQueryInfo struct {
	PageInfo
	Continent string `json:"continent"`
	Region    string `json:"region"`
	IndepYear int    `json:"indepYear"`
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

type Country struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Continent string `json:"continent"`
	Region    string `json:"region"`
	IndepYear int    `json:"indepYear"`
}

func (c *Country) TableName() string {
	return "country"
}

func (c *City) SelectPageList(p *database.Page, queryVo CityQueryInfo) error {
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

func (c *Country) SelectPageList() {

}
