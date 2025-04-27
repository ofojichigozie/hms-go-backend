package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var (
	jwtSecretKey    = []byte(os.Getenv("JWT_SECRET_KEY"))
	ErrInvalidToken = errors.New("invalid token")
)

func generateToken(userID uint, role string, duration time.Duration) (string, error) {
	claims := &JWTClaims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "iphms-go-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func GenerateTokenPair(userID uint, role string) (map[string]string, error) {
	accessToken, err := generateToken(userID, role, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(userID, role, 168*time.Hour)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}, nil
}

func VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func RefreshToken(tokenString string) (map[string]string, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokens, err := GenerateTokenPair(claims.ID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return tokens, nil
}

func GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}

func GetRoleFromToken(tokenString string) (string, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
