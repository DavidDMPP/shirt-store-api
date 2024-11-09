package repository

import (
	"shirt-store-api/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
    return r.db.Create(order).Error
}

func (r *OrderRepository) FindAll() ([]domain.Order, error) {
    var orders []domain.Order
    err := r.db.Preload("Items.Product").Preload("User").Find(&orders).Error
    return orders, err
}

func (r *OrderRepository) FindByID(id uint) (*domain.Order, error) {
    var order domain.Order
    err := r.db.Preload("Items.Product").Preload("User").First(&order, id).Error
    return &order, err
}

func (r *OrderRepository) FindByUserID(userID uint) ([]domain.Order, error) {
    var orders []domain.Order
    err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
    return orders, err
}

func (r *OrderRepository) Update(order *domain.Order) error {
    return r.db.Save(order).Error
}

func (r *OrderRepository) FindByPaymentID(paymentID string) (*domain.Order, error) {
    var order domain.Order
    err := r.db.Where("payment_id = ?", paymentID).First(&order).Error
    return &order, err
}