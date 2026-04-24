package repositories

import (
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

type ProductFilters struct {
	Search     string
	CategoryID *uint
	MinPrice   *float64
	MaxPrice   *float64
	Page       int
	PageSize   int
}

func CreateProduct(product *models.Product) error {
	return database.DB.Create(product).Error
}

func FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := database.DB.Preload("Seller").Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	loadProductStats(&product)
	return &product, nil
}

func UpdateProduct(product *models.Product) error {
	return database.DB.Save(product).Error
}

func DeleteProduct(id uint) error {
	return database.DB.Delete(&models.Product{}, id).Error
}

func ListProducts(filters ProductFilters) ([]models.Product, int64, error) {
	query := database.DB.Model(&models.Product{}).Preload("Seller").Preload("Category")

	if filters.Search != "" {
		query = query.Where("search_vector @@ plainto_tsquery('english', ?)", filters.Search)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	var total int64
	query.Count(&total)

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	offset := (filters.Page - 1) * filters.PageSize

	var products []models.Product
	err := query.Offset(offset).Limit(filters.PageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range products {
		loadProductStats(&products[i])
	}

	return products, total, nil
}

func loadProductStats(product *models.Product) {
	var result struct {
		Avg   float64
		Count int
	}
	database.DB.Raw(`
		SELECT COALESCE(AVG(rating), 0) as avg, COUNT(*) as count
		FROM reviews WHERE product_id = ?
	`, product.ID).Scan(&result)
	product.AverageRating = result.Avg
	product.ReviewCount = result.Count
}
