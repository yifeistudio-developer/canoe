package service

import "gorm.io/gorm"

type UserService struct {
	Db *gorm.DB
}
