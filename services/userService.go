package services

import (
	"errors"
	"time"

	"hells/config"
	"hells/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *models.User) error {
	db := config.GetDB()

	// Check if username already exists
	var existingUser models.User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	// Find default role if not set
	if user.RoleID == 0 {
		var defaultRole models.Role
		if err := db.Where("name = ?", "Viewer").First(&defaultRole).Error; err != nil {
			return errors.New("default role not found")
		}
		user.RoleID = defaultRole.ID
	}

	return db.Create(user).Error
}

func FindUserByEmail(email string) (*models.User, error) {
	db := config.GetDB()
	var user models.User
	err := db.Where("email = ?", email).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(userID uint) (*models.User, error) {
	db := config.GetDB()
	var user models.User
	err := db.Preload("Role").First(&user, userID).Error
	return &user, err
}

func UpdateUser(user *models.User) error {
	db := config.GetDB()
	return db.Save(user).Error
}

func ResetUserPassword(email, resetToken, newPassword string) error {
	db := config.GetDB()

	// Find and validate password reset token
	var passwordReset models.PasswordReset
	err := db.Where("email = ? AND token = ? AND expires_at > ? AND is_used = ?",
		email, resetToken, time.Now(), false).First(&passwordReset).Error

	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Begin transaction
	tx := db.Begin()

	// Update user password
	var user models.User
	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	user.PasswordHash = string(hashedPassword)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Mark reset token as used
	passwordReset.IsUsed = true
	if err := tx.Save(&passwordReset).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

func ListUsers(page, limit int) ([]models.User, int, error) {
	db := config.GetDB()
	var users []models.User
	var total int64

	// Count total users
	db.Model(&models.User{}).Count(&total)

	// Paginate and fetch users
	err := db.Preload("Role").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&users).Error

	return users, int(total), err
}
