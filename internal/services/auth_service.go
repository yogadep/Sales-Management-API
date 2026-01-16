package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"sales-management-api/internal/middlewares"
	"sales-management-api/internal/models"
	"sales-management-api/internal/repositories"
)

type AuthService struct {
	users  *repositories.UserRepo
	secret string
}

func NewAuthService(users *repositories.UserRepo, secret string) *AuthService {
	return &AuthService{users: users, secret: secret}
}

func (s *AuthService) Register(username, password string, role models.Role) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}
	return s.users.Create(&u)
}

func (s *AuthService) Login(username, password string) (string, models.Role, error) {
	u, err := s.users.FindByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	claims := middlewares.Claims{
		UserID: u.ID,
		Role:   string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", err
	}

	return token, u.Role, nil
}
