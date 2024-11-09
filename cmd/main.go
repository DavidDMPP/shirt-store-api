package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"shirt-store-api/internal/config"
	"shirt-store-api/internal/handler"
	"shirt-store-api/internal/middleware"
	"shirt-store-api/internal/repository"
	"shirt-store-api/internal/service"
	"shirt-store-api/pkg/database"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Setup Database with retry mechanism
    var db *gorm.DB
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        db, err = database.NewDatabase(cfg.Database)
        if err == nil {
            break
        }
        log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(time.Second * 5)
    }
    if err != nil {
        log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    productRepo := repository.NewProductRepository(db)
    orderRepo := repository.NewOrderRepository(db)
    settingRepo := repository.NewSettingRepository(db)

    // Initialize services
    userService := service.NewUserService(userRepo)
    productService := service.NewProductService(productRepo)
    orderService := service.NewOrderService(orderRepo, userRepo, productRepo)
    settingService := service.NewSettingService(settingRepo)
    paymentService := service.NewPaymentService(settingRepo, orderRepo)

    // Initialize handlers
    userHandler := handler.NewUserHandler(userService)
    productHandler := handler.NewProductHandler(productService)
    orderHandler := handler.NewOrderHandler(orderService, paymentService)
    settingHandler := handler.NewSettingHandler(settingService)
    paymentHandler := handler.NewPaymentHandler(paymentService, orderService)

    // Initialize Fiber
    app := fiber.New(fiber.Config{
        AppName: "Shirt Store API",
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": err.Error(),
            })
        },
    })

    // Middleware
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
        AllowMethods: "GET, POST, PUT, DELETE",
    }))

    // API routes
    api := app.Group("/api")

    // Public routes (no authentication needed)
    api.Post("/auth/register", userHandler.Register)
    api.Post("/auth/login", userHandler.Login)
    api.Post("/make-first-admin", userHandler.MakeFirstAdmin) // Public access for first setup
    api.Get("/products", productHandler.GetAllProducts)
    api.Get("/products/:id", productHandler.GetProduct)

    // Health check
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "message": "Server is healthy",
        })
    })

    // Protected routes
    protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware())

    // User routes
    protected.Get("/user/profile", userHandler.GetProfile)
    protected.Put("/user/profile", userHandler.UpdateProfile)

    // Order routes
    orders := protected.Group("/orders")
    orders.Post("/", orderHandler.CreateOrder)
    orders.Get("/", orderHandler.GetUserOrders)
    orders.Get("/:id", orderHandler.GetOrderDetail)
    orders.Post("/:id/pay", paymentHandler.ProcessPayment)
    orders.Get("/:id/status", paymentHandler.CheckPaymentStatus)

    // Admin routes
    admin := protected.Group("/admin")
    admin.Use(middleware.AdminMiddleware())
    
    // Product management (admin only)
    admin.Post("/products", productHandler.CreateProduct)
    admin.Put("/products/:id", productHandler.UpdateProduct)
    admin.Delete("/products/:id", productHandler.DeleteProduct)
    
    // Order management (admin only)
    admin.Get("/orders", orderHandler.GetAllOrders)
    admin.Put("/orders/:id/status", orderHandler.UpdateOrderStatus)

    // Settings management (admin only)
    settings := admin.Group("/settings")
    settings.Get("/midtrans", settingHandler.GetMidtransConfig)
    settings.Post("/midtrans", settingHandler.UpdateMidtransConfig)

    // Payment webhook
    api.Post("/payment/notification", paymentHandler.HandlePaymentNotification)

    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        fmt.Println("\nShutting down server...")
        _ = app.Shutdown()
    }()

    // Start server
    port := cfg.AppPort
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := app.Listen(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}