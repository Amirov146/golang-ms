package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Username  string `json:"username" gorm:"unique"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	CreatedAt time.Time
}
