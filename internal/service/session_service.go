package service

import "gorm.io/gorm"

type SessionService struct {
	Db *gorm.DB
}
