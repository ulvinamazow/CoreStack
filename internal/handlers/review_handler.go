package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/CoreStack/internal/models"
	"github.com/ulvinamazow/CoreStack/internal/repositories"
)

func CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if !repositories.UserBoughtProduct(userID.(uint), uint(productID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only review products you have purchased"})
		return
	}

	if _, err := repositories.FindReviewByUserAndProduct(userID.(uint), uint(productID)); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "you have already reviewed this product"})
		return
	}

	var req struct {
		OrderID uint   `json:"order_id" binding:"required"`
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := &models.Review{
		UserID:    userID.(uint),
		ProductID: uint(productID),
		OrderID:   req.OrderID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := repositories.CreateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func GetProductReviews(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	reviews, err := repositories.ListProductReviews(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	avg, count := repositories.GetProductAverageRating(uint(productID))

	c.JSON(http.StatusOK, gin.H{
		"reviews":        reviews,
		"average_rating": avg,
		"total_reviews":  count,
	})
}

func UpdateReview(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(*models.User)

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	review, err := repositories.FindReviewByID(uint(reviewID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if review.UserID != user.ID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Rating  *int    `json:"rating"`
		Comment *string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Rating != nil && *req.Rating >= 1 && *req.Rating <= 5 {
		review.Rating = *req.Rating
	}
	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	if err := repositories.UpdateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func DeleteReview(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(*models.User)

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	review, err := repositories.FindReviewByID(uint(reviewID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if review.UserID != user.ID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := repositories.DeleteReview(uint(reviewID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted"})
}
