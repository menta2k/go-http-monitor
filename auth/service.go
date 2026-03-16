package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Service struct {
	secret   []byte
	tokenTTL time.Duration
	users    map[string]string // username -> password (hashed in production)
}

func NewService(secret string, tokenTTL time.Duration, users map[string]string) *Service {
	return &Service{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
		users:    users,
	}
}

func (s *Service) Authenticate(username, password string) (string, error) {
	storedPass, exists := s.users[username]
	if !exists || storedPass != password {
		return "", ErrInvalidCredentials
	}

	now := time.Now().UTC()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			Issuer:    "go-http-monitor",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(tokenStr string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	return *claims, nil
}
