package repositories

import (
	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("gmail = ?", email).First(&user).Error
	return &user, err
}

func FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindUserByVerificationToken(token string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("verification_token = ?", token).First(&user).Error
	return &user, err
}

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func CountUsers() (int64, error) {
	var count int64
	err := database.DB.Model(&models.User{}).Count(&count).Error
	return count, err
}
