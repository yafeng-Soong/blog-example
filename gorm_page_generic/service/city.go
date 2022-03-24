package service

import (
	"gorm_page_generic/database"
	"gorm_page_generic/model"
)

type CityService struct{}

var cityModel model.City

// 使用泛型调用分页查询
func (c *CityService) SelectPageList(queryVo model.CityQueryInfo) (*model.PageResponse[model.City], error) {
	p := &database.Page[model.City]{
		CurrentPage: queryVo.CurrentPage,
		PageSize:    queryVo.PageSize,
	}
	wrapper := make(map[string]interface{}, 0)
	if queryVo.CountryCode != "" {
		wrapper["CountryCode"] = queryVo.CountryCode
	}
	if queryVo.District != "" {
		wrapper["District"] = queryVo.District
	}
	err := cityModel.SelectPageList(p, wrapper)
	if err != nil {
		return nil, err
	}
	pageResponse := model.NewPageResponse(p)
	return pageResponse, err
}
