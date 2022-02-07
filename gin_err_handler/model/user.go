package model

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id       int
	Username string
	Password string
	Email    string
}
