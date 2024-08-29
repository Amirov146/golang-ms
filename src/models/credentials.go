package models

type Credentials struct {
	Password string `form:"password"`
	Username string `form:"username"`
}

var Users = map[string]string{
	"user":  "user",
	"admin": "admin",
}
