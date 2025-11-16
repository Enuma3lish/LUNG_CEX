package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Enuma3lish/LUNG_CEX/backend/internal/models"
	redisClient "github.com/Enuma3lish/LUNG_CEX/backend/pkg/redis"
	"github.com/Enuma3lish/LUNG_CEX/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewPortfolioHandler(db *gorm.DB, redisClientInstance *redis.Client) *PortfolioHandler {
	return &PortfolioHandler{
		db:          db,
		redisClient: redisClientInstance,
	}
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	ctx := context.Background()

	// Try to get from cache first
	cacheKey := fmt.Sprintf("portfolio:%s", userID.String())
	if h.redisClient != nil {
		cached, err := h.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var portfolio models.PortfolioResponse
			if err := json.Unmarshal([]byte(cached), &portfolio); err == nil {
				c.JSON(http.StatusOK, portfolio)
				return
			}
		}
	}

	// Get user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get holdings
	var holdings []models.Holding
	if err := h.db.Preload("Asset").Where("user_id = ?", userID).Find(&holdings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch holdings"})
		return
	}

	// Calculate portfolio value
	totalValue := user.Balance
	holdingsWithDetails := []models.HoldingWithDetails{}

	for _, holding := range holdings {
		// Get current price (mock prices for now)
		currentPrice := utils.GetMockPrice(holding.Asset.Symbol)

		value := holding.Quantity * currentPrice
		totalValue += value

		pnl := (currentPrice - holding.AvgPrice) * holding.Quantity
		pnlPercent := ((currentPrice - holding.AvgPrice) / holding.AvgPrice) * 100

		holdingsWithDetails = append(holdingsWithDetails, models.HoldingWithDetails{
			Holding:      holding,
			CurrentPrice: currentPrice,
			Value:        value,
			PnL:          pnl,
			PnLPercent:   pnlPercent,
		})
	}

	// Calculate overall PnL
	overallPnL := totalValue - 10000.00 // Initial balance was 10000

	portfolio := models.PortfolioResponse{
		TotalValue: totalValue,
		Cash:       user.Balance,
		Holdings:   holdingsWithDetails,
		PnL:        overallPnL,
	}

	// Cache the result
	if h.redisClient != nil {
		data, _ := json.Marshal(portfolio)
		h.redisClient.Set(ctx, cacheKey, data, redisClient.PortfolioCacheTTL)
	}

	c.JSON(http.StatusOK, portfolio)
}

func (h *PortfolioHandler) GetHoldings(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	ctx := context.Background()

	// Try to get from cache first
	cacheKey := fmt.Sprintf("holdings:%s", userID.String())
	if h.redisClient != nil {
		cached, err := h.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var holdings []models.Holding
			if err := json.Unmarshal([]byte(cached), &holdings); err == nil {
				c.JSON(http.StatusOK, holdings)
				return
			}
		}
	}

	// Get from database
	var holdings []models.Holding
	if err := h.db.Preload("Asset").Where("user_id = ?", userID).Find(&holdings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch holdings"})
		return
	}

	// Cache the result
	if h.redisClient != nil {
		data, _ := json.Marshal(holdings)
		h.redisClient.Set(ctx, cacheKey, data, redisClient.HoldingsCacheTTL)
	}

	c.JSON(http.StatusOK, holdings)
}
