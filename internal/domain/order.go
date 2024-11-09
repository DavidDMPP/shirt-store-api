package domain

import (
	"time"
)

type Order struct {
    ID          uint        `json:"id" gorm:"primaryKey"`
    UserID      uint        `json:"user_id" gorm:"not null"`
    User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
    TotalAmount float64     `json:"total_amount" gorm:"not null"`
    Status      string      `json:"status" gorm:"not null;default:'pending'"`
    PaymentID   string      `json:"payment_id"`
    PaymentURL  string      `json:"payment_url"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
    ID        uint    `json:"id" gorm:"primaryKey"`
    OrderID   uint    `json:"order_id"`
    ProductID uint    `json:"product_id"`
    Product   Product `json:"product" gorm:"foreignKey:ProductID"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

type CreateOrderRequest struct {
    Items []OrderItemRequest `json:"items" validate:"required,dive"`
}

type OrderItemRequest struct {
    ProductID uint `json:"product_id" validate:"required"`
    Quantity  int  `json:"quantity" validate:"required,gt=0"`
}