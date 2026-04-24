package repositories

import (
	"time"

	"github.com/ulvinamazow/CoreStack/internal/database"
	"github.com/ulvinamazow/CoreStack/internal/models"
)

func CreateRefreshToken(token *models.RefreshToken) error {
	return database.DB.Create(token).Error
}

func FindRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := database.DB.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).First(&token).Error
	return &token, err
}

func FindRefreshTokenRaw(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := database.DB.Where("token_hash = ?", tokenHash).First(&token).Error
	return &token, err
}

func RevokeRefreshToken(tokenHash string) error {
	now := time.Now()
	return database.DB.Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", now).Error
}

func RevokeAllUserTokens(userID uint) error {
	now := time.Now()
	return database.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}
