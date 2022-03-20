package model

import (
	"gorm_page_generic/database"
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

func (c *City) SelectPageList(p *database.Page[City], wrapper map[string]interface{}) error {
	return p.SelectPage(wrapper)
}
