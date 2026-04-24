package handlers

import (
	"net/http"
	"time"

	"github.com/ulvinamazow/CoreStack/internal/config"
	"github.com/ulvinamazow/CoreStack/internal/models"
	"github.com/ulvinamazow/CoreStack/internal/repositories"
	"github.com/ulvinamazow/CoreStack/internal/services"
	"github.com/ulvinamazow/CoreStack/internal/utils"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := repositories.FindUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	if _, err := repositories.FindUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred, please try again"})
		return
	}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	count, _ := repositories.CountUsers()
	isAdmin := count == 0

	user := &models.User{
		Username:          req.Username,
		Name:              req.Name,
		Gmail:             req.Email,
		PasswordHash:      hash,
		VerificationToken: &token,
		IsAdmin:           isAdmin,
	}

	if err := repositories.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	go func() {
		_ = services.SendVerificationEmail(user.Gmail, user.Name, token)
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful, please verify your email",
		"user_id": user.ID,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repositories.FindUserByEmail(req.Email)
	if err != nil || !utils.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, req.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshTokenStr, err := utils.GenerateRandomToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	refreshExpiry := time.Duration(config.App.JWTRefreshDays) * 24 * time.Hour
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashToken(refreshTokenStr),
		ExpiresAt: time.Now().Add(refreshExpiry),
	}

	if err := repositories.CreateRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	c.SetCookie("refresh_token", refreshTokenStr, int(refreshExpiry.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshTokenStr,
		"user": gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"name":           user.Name,
			"email":          user.Gmail,
			"email_verified": user.EmailVerified,
			"is_admin":       user.IsAdmin,
		},
	})
}

func Refresh(c *gin.Context) {
	refreshTokenStr, err := c.Cookie("refresh_token")
	if err != nil {
		var req RefreshRequest
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
			return
		}
		refreshTokenStr = req.RefreshToken
	}

	tokenHash := utils.HashToken(refreshTokenStr)

	dbToken, err := repositories.FindRefreshToken(tokenHash)
	if err != nil {
		rawToken, rawErr := repositories.FindRefreshTokenRaw(tokenHash)
		if rawErr == nil && rawToken.RevokedAt != nil {
			_ = repositories.RevokeAllUserTokens(rawToken.UserID)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	if err := repositories.RevokeRefreshToken(tokenHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke token"})
		return
	}

	newRefreshStr, err := utils.GenerateRandomToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshExpiry := time.Duration(config.App.JWTRefreshDays) * 24 * time.Hour
	newRefreshToken := &models.RefreshToken{
		UserID:    dbToken.UserID,
		TokenHash: utils.HashToken(newRefreshStr),
		ExpiresAt: time.Now().Add(refreshExpiry),
	}

	if err := repositories.CreateRefreshToken(newRefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(dbToken.UserID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	c.SetCookie("refresh_token", newRefreshStr, int(refreshExpiry.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshStr,
	})
}

func Logout(c *gin.Context) {
	refreshTokenStr, _ := c.Cookie("refresh_token")
	if refreshTokenStr != "" {
		tokenHash := utils.HashToken(refreshTokenStr)
		_ = repositories.RevokeRefreshToken(tokenHash)
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}

	user, err := repositories.FindUserByVerificationToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	now := time.Now()
	user.EmailVerified = true
	user.VerificationToken = nil
	user.VerifiedAt = &now

	if err := repositories.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func ResendVerification(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(*models.User)

	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already verified"})
		return
	}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.VerificationToken = &token
	if err := repositories.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update token"})
		return
	}

	go func() {
		_ = services.SendVerificationEmail(user.Gmail, user.Name, token)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "verification email sent"})
}

func GetProfile(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(*models.User)
	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(*models.User)

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := repositories.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}
