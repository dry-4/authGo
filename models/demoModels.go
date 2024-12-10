package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserId       uint      `gorm:"unique;not null" json:"userId"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Name         string    `json:"name"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	RoleID       uint      `json:"role_id"`
	Role         Role      `gorm:"foreignkey:RoleID" json:"role"`
	Posts        []Post    `gorm:"foreignkey:UserID" json:"posts"`
	LastLogin    time.Time `gorm:"default:NULL" json:"last_login"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `gorm:"foreignkey:UserID" json:"user"`
	Status  string `gorm:"default:'draft'" json:"status"` // draft, published, archived
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
}

type PasswordReset struct {
	gorm.Model
	Email     string    `gorm:"not null" json:"email"`
	Token     string    `gorm:"not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
}
