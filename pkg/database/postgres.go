package database

import (
	"fmt"
	"log"
	"time"

	"shirt-store-api/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func NewDatabase(config *DBConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        config.Host,
        config.User,
        config.Password,
        config.DBName,
        config.Port,
        config.SSLMode,
    )

    // Disable prepared statement
    postgresConfig := postgres.Config{
        DSN: dsn,
        PreferSimpleProtocol: true, // Disable implicit prepared statement
    }

    // Custom GORM configuration
    gormConfig := &gorm.Config{
        PrepareStmt: false, // Disable prepared statement
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }

    // Connect to database
    db, err := gorm.Open(postgres.New(postgresConfig), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    // Set connection pool settings
    sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
    sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Maximum connection lifetime

    // Drop all tables first (for clean start)
    if err := db.Migrator().DropTable(
        &domain.User{},
        &domain.Product{},
        &domain.Order{},
        &domain.OrderItem{},
        &domain.MidtransConfig{},
    ); err != nil {
        log.Printf("Warning: Failed to drop tables: %v", err)
    }

    // Auto Migrate
    err = db.AutoMigrate(
        &domain.User{},
        &domain.Product{},
        &domain.Order{},
        &domain.OrderItem{},
        &domain.MidtransConfig{},
    )
    if err != nil {
        log.Printf("Failed to migrate database: %v", err)
        return nil, err
    }

    return db, nil
}

// Helper function to check database health
func CheckHealth(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}