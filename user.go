package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	Age      int    `json:"age" validate:"gte=18,lte=100"`

	Posts []Post `json:"posts"`
}

var DB *gorm.DB
