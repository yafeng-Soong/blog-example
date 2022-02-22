package model

import "gorm_page/database"

type CountryQueryInfo struct {
	PageInfo
	Continent string `json:"continent"`
	Region    string `json:"region"`
	IndepYear int    `json:"indepYear"`
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

func (c *Country) SelectPageList(p *database.Page, wrapper map[string]interface{}) error {
	err := database.SelectPage(p, wrapper, Country{})
	return err
}
