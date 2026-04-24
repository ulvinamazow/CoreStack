package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/CoreStack/internal/models"
	"github.com/ulvinamazow/CoreStack/internal/repositories"
)

func GetCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	items, err := repositories.GetCartItems(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	type CartItemResponse struct {
		models.CartItem
		DiscountedPrice float64 `json:"discounted_price"`
	}

	result := make([]CartItemResponse, len(items))
	var total float64
	for i, item := range items {
		discounted := item.Product.DiscountedPrice()
		result[i] = CartItemResponse{CartItem: item, DiscountedPrice: discounted}
		total += discounted * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{"items": result, "total": total})
}

func AddToCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := repositories.FindProductByID(req.ProductID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
		return
	}

	existing, err := repositories.FindCartItem(userID.(uint), req.ProductID)
	if err == nil {
		existing.Quantity += req.Quantity
		if err := repositories.UpdateCartItem(existing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart"})
			return
		}
		c.JSON(http.StatusOK, existing)
		return
	}

	item := &models.CartItem{
		UserID:    userID.(uint),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := repositories.CreateCartItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to cart"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func UpdateCartItem(c *gin.Context) {
	userID, _ := c.Get("user_id")

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := repositories.FindCartItemByID(uint(itemID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	item.Quantity = req.Quantity
	if err := repositories.UpdateCartItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func RemoveFromCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	if err := repositories.DeleteCartItem(uint(itemID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
