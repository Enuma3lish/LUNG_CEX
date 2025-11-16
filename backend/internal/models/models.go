package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Balance   float64   `gorm:"type:decimal(20,2);default:10000.00" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Asset represents tradeable assets
type Asset struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Symbol    string    `gorm:"unique;not null" json:"symbol"` // BTC, USDC, USDT
	Name      string    `gorm:"not null" json:"name"`
	AssetType string    `gorm:"not null" json:"asset_type"` // SPOT, FUTURES
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Holding represents user's asset holdings
type Holding struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	AssetID   uuid.UUID `gorm:"type:uuid;not null" json:"asset_id"`
	Quantity  float64   `gorm:"type:decimal(20,8);not null" json:"quantity"`
	AvgPrice  float64   `gorm:"type:decimal(20,2);not null" json:"avg_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Asset Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
}

// Trade represents a trade transaction
type Trade struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	AssetID         uuid.UUID `gorm:"type:uuid;not null" json:"asset_id"`
	TradeType       string    `gorm:"not null" json:"trade_type"` // BUY, SELL
	Quantity        float64   `gorm:"type:decimal(20,8);not null" json:"quantity"`
	Price           float64   `gorm:"type:decimal(20,2);not null" json:"price"`
	TotalAmount     float64   `gorm:"type:decimal(20,2);not null" json:"total_amount"`
	SolanaSignature string    `gorm:"type:varchar(255)" json:"solana_signature"`
	CreatedAt       time.Time `json:"created_at"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Asset Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
}

// FuturesPosition represents a futures position
type FuturesPosition struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	AssetID      uuid.UUID `gorm:"type:uuid;not null" json:"asset_id"`
	PositionType string    `gorm:"not null" json:"position_type"` // LONG, SHORT
	Quantity     float64   `gorm:"type:decimal(20,8);not null" json:"quantity"`
	EntryPrice   float64   `gorm:"type:decimal(20,2);not null" json:"entry_price"`
	Leverage     int       `gorm:"not null" json:"leverage"`
	Margin       float64   `gorm:"type:decimal(20,2);not null" json:"margin"`
	Status       string    `gorm:"not null;default:'OPEN'" json:"status"` // OPEN, CLOSED
	ClosePrice   *float64  `gorm:"type:decimal(20,2)" json:"close_price,omitempty"`
	PnL          *float64  `gorm:"type:decimal(20,2)" json:"pnl,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Asset Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
}

// DTOs for API requests/responses

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type TradeRequest struct {
	AssetSymbol string  `json:"asset_symbol" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

type FuturesTradeRequest struct {
	AssetSymbol  string  `json:"asset_symbol" binding:"required"`
	PositionType string  `json:"position_type" binding:"required,oneof=LONG SHORT"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	Leverage     int     `json:"leverage" binding:"required,min=1,max=100"`
}

type PortfolioResponse struct {
	TotalValue float64              `json:"total_value"`
	Cash       float64              `json:"cash"`
	Holdings   []HoldingWithDetails `json:"holdings"`
	PnL        float64              `json:"pnl"`
}

type HoldingWithDetails struct {
	Holding
	CurrentPrice float64 `json:"current_price"`
	Value        float64 `json:"value"`
	PnL          float64 `json:"pnl"`
	PnLPercent   float64 `json:"pnl_percent"`
}
