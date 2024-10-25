package db

import (
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
	"gorm.io/gorm"
)

type UserAdapter struct {
	db *gorm.DB
}

func (u UserAdapter) GetById(id int64) (domain.User, error) {
	return domain.User{}, nil
}

func (u UserAdapter) Save(user domain.User) error {
	return nil
}

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(64);not null;comment:'姓名';default:''"`
	Avatar    string `gorm:"type:varchar(256);not null;comment:'头像';default:''"`
	Status    int    `gorm:"type:integer;not null;comment:'状态';default:0"`
	AccountId int64  `gorm:"type:bigint;not null;comment:'关联账号id';default:0"`
}
