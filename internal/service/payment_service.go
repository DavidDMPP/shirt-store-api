package service

import (
	"encoding/json"
	"fmt"
	"time"

	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/repository"
)

type PaymentService struct {
    settingRepo *repository.SettingRepository
    orderRepo   *repository.OrderRepository
}

func NewPaymentService(settingRepo *repository.SettingRepository, orderRepo *repository.OrderRepository) *PaymentService {
    return &PaymentService{
        settingRepo: settingRepo,
        orderRepo:   orderRepo,
    }
}

type PaymentRequest struct {
    OrderID          string          `json:"order_id"`
    GrossAmount      float64         `json:"gross_amount"`
    CustomerDetails  PaymentCustomer `json:"customer_details"`
    Items           []PaymentItem    `json:"items"`
}

type PaymentCustomer struct {
    FirstName    string `json:"first_name"`
    Email       string `json:"email"`
    Phone       string `json:"phone"`
}

type PaymentItem struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Quantity    int     `json:"quantity"`
}

func (s *PaymentService) CreatePayment(order *domain.Order, user *domain.User) (string, error) {
    midtransConfig, err := s.settingRepo.GetMidtransConfig()
    if err != nil {
        return "", fmt.Errorf("failed to get Midtrans config: %v", err)
    }

    // Prepare payment request
    items := make([]PaymentItem, len(order.Items))
    for i, item := range order.Items {
        items[i] = PaymentItem{
            ID:       fmt.Sprintf("%d", item.ProductID),
            Name:     item.Product.Name,
            Price:    item.Price,
            Quantity: item.Quantity,
        }
    }

    // Create payment request
    paymentReq := &PaymentRequest{
        OrderID:     fmt.Sprintf("ORDER-%d-%d", order.ID, time.Now().Unix()),
        GrossAmount: order.TotalAmount,
        CustomerDetails: PaymentCustomer{
            FirstName: user.Name,
            Email:    user.Email,
        },
        Items: items,
    }

    // Convert to JSON for Midtrans API
    jsonData, err := json.Marshal(paymentReq)
    if err != nil {
        return "", fmt.Errorf("failed to marshal payment request: %v", err)
    }

    // Menggunakan midtransConfig untuk setup environment dan headers
    baseURL := "https://app.sandbox.midtrans.com"
    if midtransConfig.Environment == "production" {
        baseURL = "https://app.midtrans.com"
    }

    // Di sini Anda bisa menggunakan jsonData untuk request ke Midtrans
    // Contoh penggunaan dengan fmt.Printf untuk debugging
    fmt.Printf("Sending payment request to Midtrans: %s\n", string(jsonData))

    // Simulasi response dari Midtrans
    paymentURL := fmt.Sprintf("%s/snap/v1/transactions/%s?amount=%f&items=%d", 
        baseURL, 
        paymentReq.OrderID,
        paymentReq.GrossAmount,
        len(paymentReq.Items),
    )
    
    // Update order with payment information
    order.PaymentID = paymentReq.OrderID
    order.Status = "awaiting_payment"
    order.PaymentURL = paymentURL

    if err := s.orderRepo.Update(order); err != nil {
        return "", fmt.Errorf("failed to update order: %v", err)
    }

    return paymentURL, nil
}

func (s *PaymentService) HandlePaymentNotification(notification map[string]interface{}) error {
    orderID := notification["order_id"].(string)
    transactionStatus := notification["transaction_status"].(string)

    order, err := s.orderRepo.FindByPaymentID(orderID)
    if err != nil {
        return fmt.Errorf("failed to find order: %v", err)
    }

    switch transactionStatus {
    case "capture", "settlement":
        order.Status = "paid"
    case "pending":
        order.Status = "awaiting_payment"
    case "deny", "cancel", "expire":
        order.Status = "cancelled"
    default:
        order.Status = "failed"
    }

    if err := s.orderRepo.Update(order); err != nil {
        return fmt.Errorf("failed to update order status: %v", err)
    }

    return nil
}

func (s *PaymentService) GetPaymentStatus(orderID string) (string, error) {
    midtransConfig, err := s.settingRepo.GetMidtransConfig()
    if err != nil {
        return "", fmt.Errorf("failed to get Midtrans config: %v", err)
    }

    // Log untuk debugging
    fmt.Printf("Checking payment status for order %s using %s environment\n", 
        orderID, 
        midtransConfig.Environment,
    )

    // Implementasi sebenarnya akan menggunakan Midtrans SDK
    // Untuk saat ini, return status default
    return "pending", nil
}