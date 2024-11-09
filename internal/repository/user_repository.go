package repository

import (
	"shirt-store-api/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db}
}

func (r *UserRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *UserRepository) Update(user *domain.User) error {
    return r.db.Save(user).Error
}