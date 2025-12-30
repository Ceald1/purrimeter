package crypto

import "github.com/golang-jwt/jwt/v5"

// jwt type
type CustomClaims struct {
    jwt.RegisteredClaims
	Name string `json:"name"`
	IsUserToken bool `json:"is_user_token"`
}


