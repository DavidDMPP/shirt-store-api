package repository

import (
	"shirt-store-api/internal/domain"

	"gorm.io/gorm"
)

type ProductRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
    return &ProductRepository{db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
    return r.db.Create(product).Error
}

func (r *ProductRepository) FindAll() ([]domain.Product, error) {
    var products []domain.Product
    err := r.db.Find(&products).Error
    return products, err
}

func (r *ProductRepository) FindByID(id uint) (*domain.Product, error) {
    var product domain.Product
    err := r.db.First(&product, id).Error
    return &product, err
}

func (r *ProductRepository) Update(product *domain.Product) error {
    return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
    return r.db.Delete(&domain.Product{}, id).Error
}