package service

import (
	"errors"
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/repository"
)

type OrderService struct {
    orderRepo    *repository.OrderRepository
    userRepo     *repository.UserRepository
    productRepo  *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, userRepo *repository.UserRepository, productRepo *repository.ProductRepository) *OrderService {
    return &OrderService{
        orderRepo:    orderRepo,
        userRepo:     userRepo,
        productRepo:  productRepo,
    }
}

func (s *OrderService) CreateOrder(userID uint, req *domain.CreateOrderRequest) (*domain.Order, error) {
    // Verify user exists
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    // Create order items and calculate total
    var totalAmount float64
    var orderItems []domain.OrderItem

    for _, item := range req.Items {
        product, err := s.productRepo.FindByID(item.ProductID)
        if err != nil {
            return nil, errors.New("product not found")
        }

        if product.Stock < item.Quantity {
            return nil, errors.New("insufficient stock")
        }

        // Update product stock
        product.Stock -= item.Quantity
        if err := s.productRepo.Update(product); err != nil {
            return nil, err
        }

        orderItem := domain.OrderItem{
            ProductID: product.ID,
            Quantity:  item.Quantity,
            Price:     product.Price,
        }

        orderItems = append(orderItems, orderItem)
        totalAmount += product.Price * float64(item.Quantity)
    }

    // Create order
    order := &domain.Order{
        UserID:      user.ID,
        Items:       orderItems,
        TotalAmount: totalAmount,
        Status:      "pending",
    }

    if err := s.orderRepo.Create(order); err != nil {
        return nil, err
    }

    return order, nil
}

func (s *OrderService) GetUserOrders(userID uint) ([]domain.Order, error) {
    return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) GetOrderByID(id uint) (*domain.Order, error) {
    return s.orderRepo.FindByID(id)
}

func (s *OrderService) GetAllOrders() ([]domain.Order, error) {
    return s.orderRepo.FindAll()
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
    order, err := s.orderRepo.FindByID(orderID)
    if err != nil {
        return err
    }

    order.Status = status
    return s.orderRepo.Update(order)
}

func (s *OrderService) GetUserByID(userID uint) (*domain.User, error) {
    return s.userRepo.FindByID(userID)
}