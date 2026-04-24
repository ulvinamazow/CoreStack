package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/CoreStack/internal/config"
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/handlers"
	"github.com/ulvinamazow/CoreStack/internal/middleware"
)

func main() {
	config.Load()
	database.Connect()
	database.Migrate()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/refresh", handlers.Refresh)
		api.GET("/verify-email", handlers.VerifyEmail)

		auth := api.Group("")
		auth.Use(middleware.AuthRequired())
		{
			auth.POST("/logout", handlers.Logout)
			auth.POST("/resend-verification", handlers.ResendVerification)
			auth.GET("/profile", handlers.GetProfile)
			auth.PUT("/profile", handlers.UpdateProfile)
		}

		api.GET("/products", handlers.ListProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.GET("/products/:id/reviews", handlers.GetProductReviews)

		api.GET("/categories", handlers.ListCategories)

		verified := api.Group("")
		verified.Use(middleware.AuthRequired(), middleware.EmailVerifiedRequired())
		{
			verified.POST("/products", handlers.CreateProduct)
			verified.PUT("/products/:id", handlers.UpdateProduct)
			verified.DELETE("/products/:id", handlers.DeleteProduct)

			verified.GET("/cart", handlers.GetCart)
			verified.POST("/cart", handlers.AddToCart)
			verified.PUT("/cart/:item_id", handlers.UpdateCartItem)
			verified.DELETE("/cart/:item_id", handlers.RemoveFromCart)
			verified.POST("/cart/checkout", handlers.Checkout)

			verified.GET("/orders", handlers.GetOrders)
			verified.GET("/orders/:id", handlers.GetOrder)

			verified.POST("/products/:id/reviews", handlers.CreateReview)
		}

		reviews := api.Group("")
		reviews.Use(middleware.AuthRequired())
		{
			reviews.PUT("/reviews/:id", handlers.UpdateReview)
			reviews.DELETE("/reviews/:id", handlers.DeleteReview)
		}

		admin := api.Group("")
		admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
		{
			admin.POST("/categories", handlers.CreateCategory)
		}
	}

	r.POST("/api/webhooks/stripe", handlers.StripeWebhook)

	log.Printf("Server starting on port %s", config.App.Port)
	if err := r.Run(":" + config.App.Port); err != nil {
		log.Fatal(err)
	}
}
