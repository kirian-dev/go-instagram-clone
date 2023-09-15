package utils

import (
	"errors"
	"fmt"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}
type RefreshToken struct {
	Token     string
	ExpiresAt time.Time
	Email     string
	Role      string
}

func GenerateJWTTokens(email, role string, c *config.Config) (string, string, error) {
	// Access Token
	accessSigningKey := []byte(c.JwtSecretKey)
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["authorized"] = true
	accessClaims["email"] = email
	accessClaims["role"] = role
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessTokenString, err := accessToken.SignedString(accessSigningKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshSigningKey := []byte(c.JwtSecretKey)
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["authorized"] = true
	refreshClaims["email"] = email
	refreshClaims["role"] = role
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	refreshTokenString, err := refreshToken.SignedString(refreshSigningKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateJWTToken(tokenString string, cfg *config.Config) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(e.ErrInvalidToken)
}
