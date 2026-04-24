package repositories

import (
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

func CreateCategory(category *models.Category) error {
	return database.DB.Create(category).Error
}

func FindCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := database.DB.First(&category, id).Error
	return &category, err
}

func ListCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Where("parent_id IS NULL").Preload("Children").Preload("Children.Children").Find(&categories).Error
	return categories, err
}

func UpdateCategory(category *models.Category) error {
	return database.DB.Save(category).Error
}

func DeleteCategory(id uint) error {
	return database.DB.Delete(&models.Category{}, id).Error
}
