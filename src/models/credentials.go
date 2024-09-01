package models

type Credentials struct {
	Password string `form:"password"`
	Username string `form:"username"`
}
