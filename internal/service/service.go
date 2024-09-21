package service

import (
	"github.com/kataras/golog"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupServices(d *gorm.DB, l *golog.Logger) {
	if d == nil {
		panic("database is nil")
	}
	db = d
	logger = l
}
