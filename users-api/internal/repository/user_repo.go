package repository

import (
	"errors"

	"users-api/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}

type userRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) UserRepository { return &userRepo{db: db} }
func AutoMigrate(db *gorm.DB) error          { return db.AutoMigrate(&domain.User{}) }

func (r *userRepo) Create(u *domain.User) error { return r.db.Create(u).Error }

func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}
