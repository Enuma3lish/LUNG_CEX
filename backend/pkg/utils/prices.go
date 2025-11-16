package utils

import (
	"math/rand"
	"time"
)

// Mock prices for virtual trading
var basePrices = map[string]float64{
	"BTC":      45000.00,
	"ETH":      2500.00,
	"SOL":      100.00,
	"USDC":     1.00,
	"USDT":     1.00,
	"BTC-PERP": 45000.00,
	"ETH-PERP": 2500.00,
	"SOL-PERP": 100.00,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetMockPrice returns a mock price with slight random variation
func GetMockPrice(symbol string) float64 {
	basePrice, exists := basePrices[symbol]
	if !exists {
		return 1.00
	}

	// Add random variation (-2% to +2%) for non-stablecoins
	if symbol != "USDC" && symbol != "USDT" {
		variation := (rand.Float64() - 0.5) * 0.04 // -2% to +2%
		return basePrice * (1 + variation)
	}

	return basePrice
}

// GetAllPrices returns all current mock prices
func GetAllPrices() map[string]float64 {
	prices := make(map[string]float64)
	for symbol := range basePrices {
		prices[symbol] = GetMockPrice(symbol)
	}
	return prices
}
