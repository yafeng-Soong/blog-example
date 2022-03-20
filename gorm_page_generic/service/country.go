package service

import (
	"gorm_page_generic/database"
	"gorm_page_generic/model"
)

type CountryService struct{}

var countryModel model.Country

func (c *CountryService) SelectPageList(queryVo model.CountryQueryInfo) (*model.PageResponse[model.Country], error) {
	p := &database.Page[model.Country]{
		CurrentPage: queryVo.CurrentPage,
		PageSize:    queryVo.PageSize,
	}
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
	err := countryModel.SelectPageList(p, wrapper)
	if err != nil {
		return nil, err
	}
	pageResponse := model.NewPageResponse(p)
	return pageResponse, nil
}
