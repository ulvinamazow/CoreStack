package repositories

import (
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

func CreateReview(review *models.Review) error {
	return database.DB.Create(review).Error
}

func FindReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	err := database.DB.Preload("User").First(&review, id).Error
	return &review, err
}

func FindReviewByUserAndProduct(userID, productID uint) (*models.Review, error) {
	var review models.Review
	err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	return &review, err
}

func ListProductReviews(productID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := database.DB.Preload("User").Where("product_id = ?", productID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func UpdateReview(review *models.Review) error {
	return database.DB.Save(review).Error
}

func DeleteReview(id uint) error {
	return database.DB.Delete(&models.Review{}, id).Error
}

func GetProductAverageRating(productID uint) (float64, int64) {
	var result struct {
		Avg   float64
		Count int64
	}
	database.DB.Raw(`
		SELECT COALESCE(AVG(rating), 0) as avg, COUNT(*) as count
		FROM reviews WHERE product_id = ?
	`, productID).Scan(&result)
	return result.Avg, result.Count
}
