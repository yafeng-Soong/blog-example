package service

import (
	"gin_err_handler/errors"
	"gin_err_handler/model"
)

type UserService struct{}

var db = model.User{
	Id:       1,
	Username: "Alice",
	Email:    "alice@gmail.cn",
	Password: "123456",
} // 将用户数据写死在代码里

func (s *UserService) Login(login model.Login) (*model.User, error) {
	u := &model.User{}
	if login.Email != db.Email {
		return nil, errors.LOGIN_UNKNOWN
	}
	if login.Password != db.Password {
		return nil, errors.LOGIN_ERROR
	}
	*u = db
	u.Password = "" // 密码是敏感信息不返回
	return u, nil
}
