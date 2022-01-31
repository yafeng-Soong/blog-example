package model

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id       int
	UserName string
	Email    string
}

func NewInfo(user User) *UserInfo {
	return &UserInfo{
		Id:       user.Id,
		UserName: user.UserName,
		Email:    user.Email,
	}
}
