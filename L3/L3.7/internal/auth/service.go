package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service - сервис для работы с JWT токенами
type Service struct {
	secret string
}

// NewService - создание сервиса авторизации
func NewService(secret string) *Service {
	return &Service{
		secret: secret,
	}
}

// GenerateToken - генерация JWT токена
func (s *Service) GenerateToken(username string, role string) (string, error) {

	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}

// ParseToken - разбор JWT токена и получение данных пользователя
func (s *Service) ParseToken(tokenStr string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
