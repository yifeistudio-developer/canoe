package service

import (
	"github.com/kataras/golog"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

}

func SetDB(d *gorm.DB) {
	if d == nil {
		panic("database is nil")
	}
	db = d
}

func SetLogger(d *golog.Logger) {
	logger = d
}
