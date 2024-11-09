package service

import (
	"errors"
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/repository"
	"shirt-store-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo}
}

func (s *UserService) Register(req *domain.RegisterRequest) error {
    // Check if email already exists
    _, err := s.repo.FindByEmail(req.Email)
    if err == nil {
        return errors.New("email already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &domain.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     "user", // Default role
    }

    return s.repo.Create(user)
}

func (s *UserService) Login(req *domain.LoginRequest) (string, error) {
    user, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    // Compare passwords
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token
    token, err := jwt.GenerateToken(user.ID, user.Role)
    if err != nil {
        return "", err
    }

    return token, nil
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
    return s.repo.FindByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
    return s.repo.FindByEmail(email)
}

func (s *UserService) UpdateUser(user *domain.User) error {
    return s.repo.Update(user)
}