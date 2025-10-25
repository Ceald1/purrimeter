package db

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)


type AgentClaims struct {
	Username string `json:"username"`
    jwt.RegisteredClaims
}


func CreateToken(agentName string) (string, error) {
	secret := []byte(os.Getenv("SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AgentClaims{Username: agentName})
	return token.SignedString(secret)
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
    secretKey := []byte(os.Getenv("SECRET"))
    
    // Parse and verify the token
    token, err := jwt.ParseWithClaims(tokenString, &AgentClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify the signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    // Extract and validate claims
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}


func DecodeWithoutKey(tokenString string) (string, error) {
    token, _, err := jwt.NewParser().ParseUnverified(tokenString, &AgentClaims{})
    if err != nil {
		return "", err
	}
    if claims, ok := token.Claims.(*AgentClaims); ok {
		return claims.Username, nil
	}
    return "", err
}