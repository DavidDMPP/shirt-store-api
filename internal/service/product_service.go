package service

import (
	"errors"
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/repository"
)

type ProductService struct {
    repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
    return &ProductService{repo}
}

func (s *ProductService) CreateProduct(req *domain.CreateProductRequest) error {
    product := &domain.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        ImageURL:    req.ImageURL,
    }

    return s.repo.Create(product)
}

func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
    return s.repo.FindAll()
}

func (s *ProductService) GetProductByID(id uint) (*domain.Product, error) {
    return s.repo.FindByID(id)
}

func (s *ProductService) UpdateProduct(id uint, req *domain.CreateProductRequest) error {
    product, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }

    product.Name = req.Name
    product.Description = req.Description
    product.Price = req.Price
    product.Stock = req.Stock
    product.ImageURL = req.ImageURL

    return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
    return s.repo.Delete(id)
}

func (s *ProductService) UpdateStock(id uint, quantity int) error {
    product, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }

    if product.Stock < quantity {
        return errors.New("insufficient stock")
    }

    product.Stock -= quantity
    return s.repo.Update(product)
}