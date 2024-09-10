package service

import "gorm.io/gorm"

var db *gorm.DB

var Us *UserService

func init() {

}

func SetDB(d *gorm.DB) {
	if d == nil {
		panic("database is nil")
	}
	db = d
}
