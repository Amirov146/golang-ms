package models

import (
	"github.com/vk-rv/pvx"
	"time"
)

type AdditionalClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Footer struct {
	MetaData string `json:"meta_data"`
}

type ServiceClaims struct {
	pvx.RegisteredClaims
	AdditionalClaims
	Footer
}

type TokenData struct {
	Subject  string
	Duration time.Duration
	AdditionalClaims
	Footer
}
