package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "random_salt_string"
	signingKey = "super_secret_key"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GetUser(username string) (entity.User, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *AuthService) CreateUser(username, password string) error {
	hashedPassword := generatePasswordHash(password)

	newUser := entity.User{
		Username: username,
		Password: hashedPassword,
		Coins:    1000,
	}

	if _, err := s.repo.CreateUser(newUser); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	hashedPassword := generatePasswordHash(password)
	if user.Password != hashedPassword {
		return "", errors.New("incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", hash)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}
