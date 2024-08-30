package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primary_key;autoIncrement"`
	Username  string `json:"username" gorm:"unique" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" gorm:"unique" binding:"required"`
	CreatedAt time.Time
	//Roles     []Role `gorm:"many2many:users_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID"`
	Roles []Role `gorm:"many2many:users_roles;"`
}
