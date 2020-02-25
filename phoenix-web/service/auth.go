package service

import "errors"

type LoginService struct {
	UserName string `form:"username" json:"username" binding:"required,min=5,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=20"`
}

func (s *LoginService) Login() error {
	if s.UserName == "admin" && s.Password == "admin" {
		return nil
	}
	return errors.New("Login Failed")
}
