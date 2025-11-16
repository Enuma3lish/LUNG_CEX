package main

import (
	"log"
	"os"

	"github.com/Enuma3lish/LUNG_CEX/backend/internal/database"
	"github.com/Enuma3lish/LUNG_CEX/backend/internal/handlers"
	"github.com/Enuma3lish/LUNG_CEX/backend/internal/middleware"
	"github.com/Enuma3lish/LUNG_CEX/backend/pkg/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Redis
	redisClient := redis.InitRedis()

	// Initialize Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	tradeHandler := handlers.NewTradeHandler(db, redisClient)
	portfolioHandler := handlers.NewPortfolioHandler(db, redisClient)

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Trading endpoints
		protected.POST("/trade/buy", tradeHandler.BuyAsset)
		protected.POST("/trade/sell", tradeHandler.SellAsset)
		protected.GET("/trades/history", tradeHandler.GetTradeHistory)

		// Portfolio endpoints
		protected.GET("/portfolio", portfolioHandler.GetPortfolio)
		protected.GET("/portfolio/holdings", portfolioHandler.GetHoldings)

		// User endpoints
		protected.GET("/user/profile", authHandler.GetProfile)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
