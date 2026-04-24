package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/ulvinamazow/CoreStack/internal/config"
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
	"github.com/ulvinamazow/CoreStack/internal/repositories"
)

func Checkout(c *gin.Context) {
	userID, _ := c.Get("user_id")

	items, err := repositories.GetCartItems(userID.(uint))
	if err != nil || len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	var total float64
	var orderItems []models.OrderItem

	for _, item := range items {
		product, err := repositories.FindProductByID(item.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %d not found", item.ProductID)})
			return
		}
		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("insufficient stock for %s", product.Name)})
			return
		}
		unitPrice := product.DiscountedPrice()
		total += unitPrice * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			UnitPrice:       unitPrice,
			DiscountApplied: product.Discount,
		})
	}

	stripe.Key = config.App.StripeSecretKey
	amountCents := int64(total * 100)

	intentParams := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(config.App.StripeCurrency),
	}
	intentParams.AddMetadata("user_id", fmt.Sprintf("%d", userID.(uint)))

	pi, err := paymentintent.New(intentParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment intent"})
		return
	}

	piID := pi.ID
	order := &models.Order{
		UserID:                userID.(uint),
		TotalAmount:           total,
		StripePaymentIntentID: &piID,
		Status:                "pending",
	}

	tx := database.DB.Begin()
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order items"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"order_id":      order.ID,
		"client_secret": pi.ClientSecret,
		"total":         total,
	})
}

func GetOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	orders, err := repositories.ListUserOrders(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := repositories.FindOrderByUserAndID(userID.(uint), uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
