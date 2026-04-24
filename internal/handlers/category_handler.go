package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/CoreStack/internal/models"
	"github.com/ulvinamazow/CoreStack/internal/repositories"
)

func ListCategories(c *gin.Context) {
	categories, err := repositories.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required,min=1,max=100"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if err := repositories.CreateCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category updated"})
}
