package config

import (
	"fmt"
	"os"
	"path/filepath"
	"shirt-store-api/pkg/database"

	"github.com/joho/godotenv"
)

type Config struct {
    AppPort      string
    AppEnv       string
    Database     *database.DBConfig
    JWTSecret    string
}

func LoadConfig() (*Config, error) {
    // Get executable directory
    ex, err := os.Executable()
    if err != nil {
        return nil, fmt.Errorf("failed to get executable path: %v", err)
    }
    exPath := filepath.Dir(ex)

    // Possible paths for .env file
    paths := []string{
        filepath.Join(exPath, ".env"),
        ".env",
        "../.env",
        "../../.env",
    }

    // Try to load .env from possible paths
    var loaded bool
    for _, path := range paths {
        absPath, _ := filepath.Abs(path)
        if _, err := os.Stat(absPath); err == nil {
            err = godotenv.Load(absPath)
            if err == nil {
                loaded = true
                fmt.Printf("Loaded .env from: %s\n", absPath)
                break
            }
        }
    }

    if !loaded && os.Getenv("APP_ENV") != "production" {
        return nil, fmt.Errorf("no .env file found in any of the possible locations")
    }

    dbConfig := &database.DBConfig{
        Host:     getEnv("SUPABASE_HOST", ""),
        Port:     getEnv("SUPABASE_PORT", "6543"),
        User:     getEnv("SUPABASE_USER", ""),
        Password: getEnv("SUPABASE_PASSWORD", ""),
        DBName:   getEnv("SUPABASE_DB_NAME", "postgres"),
        SSLMode:  getEnv("SUPABASE_SSL_MODE", "require"),
    }

    return &Config{
        AppPort:   getEnv("PORT", "8080"),
        AppEnv:    getEnv("APP_ENV", "development"),
        Database:  dbConfig,
        JWTSecret: getEnv("JWT_SECRET", "0a2e0aad-8bc4-4f20-a628-027c15055143"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}