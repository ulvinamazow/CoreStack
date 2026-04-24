package repositories

import (
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

func CreateOrder(order *models.Order) error {
	return database.DB.Create(order).Error
}

func FindOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := database.DB.Preload("Items").Preload("Items.Product").First(&order, id).Error
	return &order, err
}

func FindOrderByUserAndID(userID, orderID uint) (*models.Order, error) {
	var order models.Order
	err := database.DB.Preload("Items").Preload("Items.Product").
		Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	return &order, err
}

func ListUserOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := database.DB.Preload("Items").Where("user_id = ?", userID).
		Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func FindOrderByPaymentIntent(intentID string) (*models.Order, error) {
	var order models.Order
	err := database.DB.Where("stripe_payment_intent_id = ?", intentID).First(&order).Error
	return &order, err
}

func UpdateOrder(order *models.Order) error {
	return database.DB.Save(order).Error
}

func CreateOrderItems(items []models.OrderItem) error {
	return database.DB.Create(&items).Error
}

func UserBoughtProduct(userID, productID uint) bool {
	var count int64
	database.DB.Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = 'paid'", userID, productID).
		Count(&count)
	return count > 0
}
