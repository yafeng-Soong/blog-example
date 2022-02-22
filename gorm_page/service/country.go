package service

import (
	"gorm_page/database"
	"gorm_page/model"
)

type CountryService struct{}

var countryModel model.Country

func (c *CountryService) SelectPageList(p *database.Page, queryVo model.CountryQueryInfo) error {
	p.CurrentPage = queryVo.CurrentPage
	p.PageSize = queryVo.PageSize
	wrapper := make(map[string]interface{}, 0)
	if queryVo.Continent != "" {
		wrapper["Continent"] = queryVo.Continent
	}
	if queryVo.Region != "" {
		wrapper["Region"] = queryVo.Region
	}
	if queryVo.IndepYear != 0 {
		wrapper["IndepYear"] = queryVo.IndepYear
	}
	return countryModel.SelectPageList(p, wrapper)
}
