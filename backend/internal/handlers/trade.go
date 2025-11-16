package handlers

import (
	"fmt"
	"net/http"

	"github.com/Enuma3lish/LUNG_CEX/backend/internal/models"
	"github.com/Enuma3lish/LUNG_CEX/backend/pkg/blockchain"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TradeHandler struct {
	db            *gorm.DB
	redisClient   *redis.Client
	solanaClient  *blockchain.SolanaClient
}

func NewTradeHandler(db *gorm.DB, redisClient *redis.Client) *TradeHandler {
	solanaClient, err := blockchain.NewSolanaClient()
	if err != nil {
		// Log error but continue - blockchain integration is optional
		fmt.Printf("Warning: Failed to initialize Solana client: %v\n", err)
	}

	return &TradeHandler{
		db:           db,
		redisClient:  redisClient,
		solanaClient: solanaClient,
	}
}

func (h *TradeHandler) BuyAsset(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find asset
	var asset models.Asset
	if err := tx.Where("symbol = ?", req.AssetSymbol).First(&asset).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	// Get user
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Calculate total cost
	totalCost := req.Quantity * req.Price

	// Check if user has enough balance
	if user.Balance < totalCost {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Deduct balance
	user.Balance -= totalCost
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Record trade on blockchain (async, non-blocking)
	var solanaSignature string
	if h.solanaClient != nil {
		sig, err := h.solanaClient.RecordTradeOnChain(
			userID,
			req.AssetSymbol,
			"BUY",
			req.Quantity,
			req.Price,
		)
		if err == nil {
			solanaSignature = sig
		}
	}

	// Create trade record
	trade := models.Trade{
		UserID:          userID,
		AssetID:         asset.ID,
		TradeType:       "BUY",
		Quantity:        req.Quantity,
		Price:           req.Price,
		TotalAmount:     totalCost,
		SolanaSignature: solanaSignature,
	}

	if err := tx.Create(&trade).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trade"})
		return
	}

	// Update or create holding
	var holding models.Holding
	result := tx.Where("user_id = ? AND asset_id = ?", userID, asset.ID).First(&holding)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new holding
		holding = models.Holding{
			UserID:   userID,
			AssetID:  asset.ID,
			Quantity: req.Quantity,
			AvgPrice: req.Price,
		}
		if err := tx.Create(&holding).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create holding"})
			return
		}
	} else {
		// Update existing holding (calculate new average price)
		totalQuantity := holding.Quantity + req.Quantity
		holding.AvgPrice = ((holding.AvgPrice * holding.Quantity) + (req.Price * req.Quantity)) / totalQuantity
		holding.Quantity = totalQuantity

		if err := tx.Save(&holding).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update holding"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Invalidate cache
	if h.redisClient != nil {
		h.redisClient.Del(c, fmt.Sprintf("portfolio:%s", userID.String()))
		h.redisClient.Del(c, fmt.Sprintf("holdings:%s", userID.String()))
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Trade executed successfully",
		"trade":             trade,
		"solana_signature":  solanaSignature,
		"remaining_balance": user.Balance,
	})
}

func (h *TradeHandler) SellAsset(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find asset
	var asset models.Asset
	if err := tx.Where("symbol = ?", req.AssetSymbol).First(&asset).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	// Get user
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get holding
	var holding models.Holding
	if err := tx.Where("user_id = ? AND asset_id = ?", userID, asset.ID).First(&holding).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No holding found for this asset"})
		return
	}

	// Check if user has enough quantity
	if holding.Quantity < req.Quantity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient quantity"})
		return
	}

	// Calculate total proceeds
	totalProceeds := req.Quantity * req.Price

	// Add to balance
	user.Balance += totalProceeds
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Record trade on blockchain
	var solanaSignature string
	if h.solanaClient != nil {
		sig, err := h.solanaClient.RecordTradeOnChain(
			userID,
			req.AssetSymbol,
			"SELL",
			req.Quantity,
			req.Price,
		)
		if err == nil {
			solanaSignature = sig
		}
	}

	// Create trade record
	trade := models.Trade{
		UserID:          userID,
		AssetID:         asset.ID,
		TradeType:       "SELL",
		Quantity:        req.Quantity,
		Price:           req.Price,
		TotalAmount:     totalProceeds,
		SolanaSignature: solanaSignature,
	}

	if err := tx.Create(&trade).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trade"})
		return
	}

	// Update holding
	holding.Quantity -= req.Quantity

	if holding.Quantity == 0 {
		// Delete holding if quantity is zero
		if err := tx.Delete(&holding).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete holding"})
			return
		}
	} else {
		if err := tx.Save(&holding).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update holding"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Invalidate cache
	if h.redisClient != nil {
		h.redisClient.Del(c, fmt.Sprintf("portfolio:%s", userID.String()))
		h.redisClient.Del(c, fmt.Sprintf("holdings:%s", userID.String()))
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Trade executed successfully",
		"trade":             trade,
		"solana_signature":  solanaSignature,
		"remaining_balance": user.Balance,
	})
}

func (h *TradeHandler) GetTradeHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var trades []models.Trade
	if err := h.db.Preload("Asset").Where("user_id = ?", userID).Order("created_at DESC").Limit(100).Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trade history"})
		return
	}

	c.JSON(http.StatusOK, trades)
}
