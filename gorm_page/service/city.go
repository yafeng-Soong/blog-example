package service

import (
	"gorm_page/database"
	"gorm_page/model"
)

type CityService struct{}

var cityModel model.City

// 使用反射调用分页查询
func (c *CityService) SelectPageList(p *database.Page, queryVo model.CityQueryInfo) error {
	p.CurrentPage = queryVo.CurrentPage
	p.PageSize = queryVo.PageSize
	wrapper := make(map[string]interface{}, 0)
	if queryVo.CountryCode != "" {
		wrapper["CountryCode"] = queryVo.CountryCode
	}
	if queryVo.District != "" {
		wrapper["District"] = queryVo.District
	}
	err := cityModel.SelectPageList(p, wrapper)
	return err
}

// 不适用反射调用分页查询
func (c *CityService) SelectPageList1(p *database.Page, queryVo model.CityQueryInfo) error {
	p.CurrentPage = queryVo.CurrentPage
	p.PageSize = queryVo.PageSize
	wrapper := make(map[string]interface{}, 0)
	if queryVo.CountryCode != "" {
		wrapper["CountryCode"] = queryVo.CountryCode
	}
	if queryVo.District != "" {
		wrapper["District"] = queryVo.District
	}
	p.Total = cityModel.CountAll(wrapper)
	if p.Total == 0 {
		return nil // 若记录总数为0直接返回，不再执行Limit查询
	}
	return cityModel.SelectList(p, wrapper)
}
