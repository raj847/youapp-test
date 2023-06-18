package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(50);not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:text;not null"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	Username string `json:"username" gorm:"type:varchar(50);not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:text;not null"`
}
