package models

type UsersRolesID struct {
	UserID uint `gorm:"column:user_id;primaryKey"`
	RoleID uint `gorm:"column:role_id;primaryKey"`
}

type UsersRoles struct {
	UsersRolesID UsersRolesID `gorm:"embedded;embeddedPrefix:"`
}
