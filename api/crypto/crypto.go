package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// verifies the token and returns the `CustomClaims` struct and error
func VerifyToken(jwtToken, secret string) (claims *CustomClaims, err error) {
	claims = &CustomClaims{} // Initialize claims pointer
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil // Convert to []byte
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}


// --- cryptography and JWT tokens -----

func Hash(input string) (out string) {
	h := sha256.New()
	h.Write([]byte(input))
	hashed := h.Sum(nil)
	out = hex.EncodeToString(hashed)
	return out
}
func HashToNumber(input string) *big.Int {
    // Generate SHA-256 hash
    hash := sha256.Sum256([]byte(input))
    
    // Convert to big.Int
    return new(big.Int).SetBytes(hash[:])
}

// create a jwt token using username or name, if token is a user set to `True` if agent `False`
func CreateToken(username string, isUser bool) (jwt_token string, err error) {
	secret := os.Getenv("JWT_SECRET")
	claims := CustomClaims{
		Name: username,
		IsUserToken: isUser,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

