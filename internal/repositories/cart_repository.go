package repositories

import (
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

func GetCartItems(userID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := database.DB.Preload("Product").Preload("Product.Category").
		Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func FindCartItem(userID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	return &item, err
}

func FindCartItemByID(id, userID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	return &item, err
}

func CreateCartItem(item *models.CartItem) error {
	return database.DB.Create(item).Error
}

func UpdateCartItem(item *models.CartItem) error {
	return database.DB.Save(item).Error
}

func DeleteCartItem(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{}).Error
}

func ClearCart(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}
