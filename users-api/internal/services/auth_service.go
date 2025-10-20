package services

import (
	"errors"
	"time"

	"users-api/internal/domain"
	"users-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterNormal(name, email, password string) (*domain.User, error)
	Login(email, password string) (token string, exp time.Time, user *domain.User, err error)
}

type authService struct {
	users     repository.UserRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthService(users repository.UserRepository, secret string, ttl time.Duration) AuthService {
	return &authService{users: users, jwtSecret: []byte(secret), jwtTTL: ttl}
}

func (s *authService) RegisterNormal(name, email, password string) (*domain.User, error) {
	existing, err := s.users.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email_already_in_use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{Name: name, Email: email, Password: string(hash), Role: domain.RoleNormal}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	u.Password = ""
	return u, nil
}

func (s *authService) Login(email, password string) (string, time.Time, *domain.User, error) {
	u, err := s.users.FindByEmail(email)
	if err != nil {
		return "", time.Time{}, nil, err
	}
	if u == nil {
		return "", time.Time{}, nil, errors.New("invalid_credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", time.Time{}, nil, errors.New("invalid_credentials")
	}

	now := time.Now()
	exp := now.Add(s.jwtTTL)
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"role": string(u.Role),
		"iat":  now.Unix(),
		"exp":  exp.Unix(),
		"iss":  "users-api",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, nil, err
	}

	u.Password = ""
	return token, exp, u, nil
}
