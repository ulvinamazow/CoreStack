package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/ulvinamazow/CoreStack/internal/config"
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/repositories"
)

func StripeWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), config.App.StripeWebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook signature"})
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		handlePaymentSucceeded(event)
	case "payment_intent.payment_failed":
		handlePaymentFailed(event)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func handlePaymentSucceeded(event stripe.Event) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		log.Printf("failed to parse payment intent: %v", err)
		return
	}

	order, err := repositories.FindOrderByPaymentIntent(pi.ID)
	if err != nil {
		log.Printf("order not found for payment intent %s: %v", pi.ID, err)
		return
	}

	tx := database.DB.Begin()

	order.Status = "paid"
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		log.Printf("failed to update order status: %v", err)
		return
	}

	fullOrder, err := repositories.FindOrderByID(order.ID)
	if err != nil {
		tx.Rollback()
		log.Printf("failed to load order items: %v", err)
		return
	}

	for _, item := range fullOrder.Items {
		if err := tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?",
			item.Quantity, item.ProductID, item.Quantity).Error; err != nil {
			tx.Rollback()
			log.Printf("failed to update stock: %v", err)
			return
		}
	}

	if err := tx.Where("user_id = ?", order.UserID).Delete(&struct {
		UserID uint
	}{}).Error; err != nil {
		tx.Rollback()
		log.Printf("failed to clear cart: %v", err)
		return
	}

	if err := tx.Exec("DELETE FROM cart_items WHERE user_id = ?", order.UserID).Error; err != nil {
		tx.Rollback()
		log.Printf("failed to clear cart: %v", err)
		return
	}

	tx.Commit()
	log.Printf("order %d marked as paid", order.ID)
}

func handlePaymentFailed(event stripe.Event) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		log.Printf("failed to parse payment intent: %v", err)
		return
	}

	order, err := repositories.FindOrderByPaymentIntent(pi.ID)
	if err != nil {
		log.Printf("order not found for payment intent %s", pi.ID)
		return
	}

	order.Status = "cancelled"
	if err := repositories.UpdateOrder(order); err != nil {
		log.Printf("failed to update order status: %v", err)
	}
}
