package entity

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	Displayname string    `json:"display_name" gorm:"type:varchar(255)"`
	Gender      string    `json:"gender"`
	Birthday    time.Time `json:"birthday" `
	Horoscope   string    `json:"horoscope" gorm:"type:varchar(255)"`
	Zodiac      string    `json:"zodiac" gorm:"type:varchar(255)"`
	Height      int       `json:"height" `
	Weight      int       `json:"weight" `
	PhotoURL1   string    `json:"foto" gorm:"type:varchar(255)"`
}

type ProfileReq struct {
	Displayname string `json:"display_name" gorm:"type:varchar(255)"`
	Gender      string `json:"gender"`
	Birthday    string `json:"birthday" `
	Horoscope   string `json:"horoscope" gorm:"type:varchar(255)"`
	Zodiac      string `json:"zodiac" gorm:"type:varchar(255)"`
	Height      int    `json:"height" `
	Weight      int    `json:"weight" `
	PhotoURL1   string `json:"foto" gorm:"type:varchar(255)"`
}
