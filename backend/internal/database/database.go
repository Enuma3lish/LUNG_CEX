package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Enuma3lish/LUNG_CEX/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "lung_cex"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Asset{},
		&models.Holding{},
		&models.Trade{},
		&models.FuturesPosition{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed initial assets
	if err := seedAssets(db); err != nil {
		return fmt.Errorf("failed to seed assets: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

func seedAssets(db *gorm.DB) error {
	assets := []models.Asset{
		{Symbol: "USDC", Name: "USD Coin", AssetType: "SPOT"},
		{Symbol: "USDT", Name: "Tether USD", AssetType: "SPOT"},
		{Symbol: "BTC", Name: "Bitcoin", AssetType: "SPOT"},
		{Symbol: "ETH", Name: "Ethereum", AssetType: "SPOT"},
		{Symbol: "SOL", Name: "Solana", AssetType: "SPOT"},
		{Symbol: "BTC-PERP", Name: "Bitcoin Perpetual Futures", AssetType: "FUTURES"},
		{Symbol: "ETH-PERP", Name: "Ethereum Perpetual Futures", AssetType: "FUTURES"},
		{Symbol: "SOL-PERP", Name: "Solana Perpetual Futures", AssetType: "FUTURES"},
	}

	for _, asset := range assets {
		var existing models.Asset
		result := db.Where("symbol = ?", asset.Symbol).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&asset).Error; err != nil {
				return err
			}
			log.Printf("Created asset: %s", asset.Symbol)
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
