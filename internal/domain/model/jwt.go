package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Id           int32
	Username     string
	TenantId     int32
	IsSuperAdmin bool

	jwt.RegisteredClaims
}
