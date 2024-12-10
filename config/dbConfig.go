package config

import (
	"fmt"
	"log"

	"hells/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConfig holds db configrations
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Database connection instance
var dbConfig = DBConfig{
	Host:     "127.0.0.1",
	Port:     "3306",
	User:     "root",
	Password: "Sudo@123",
	DBName:   "store",
}

func InitDatabase() (*gorm.DB, error) {
	// Database connection parameters
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := InitializeRoles(db); err != nil {
		return nil, fmt.Errorf("failed to initialize roles: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Post{},
		&models.Permission{},
		&models.PasswordReset{},
	)
	if err != nil {
		return nil, fmt.Errorf("database migration failed: %v", err)
	}

	return db, nil
}

func SeedDatabase(db *gorm.DB) error {
	// Create default roles
	roles := []models.Role{
		{Name: "Admin", Description: "Full platform access"},
		{Name: "Editor", Description: "Can create and manage posts"},
		{Name: "Viewer", Description: "Read-only access"},
	}

	for _, role := range roles {
		var existingRole models.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		if result.Error != nil {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Error creating role %s: %v", role.Name, err)
			}
		}
	}

	// Create default permissions
	permissions := []models.Permission{
		{Name: "create_post", Description: "Create new blog posts"},
		{Name: "edit_post", Description: "Edit existing blog posts"},
		{Name: "delete_post", Description: "Delete blog posts"},
		{Name: "manage_users", Description: "Manage user accounts"},
	}

	for _, perm := range permissions {
		var existingPerm models.Permission
		result := db.Where("name = ?", perm.Name).First(&existingPerm)
		if result.Error != nil {
			if err := db.Create(&perm).Error; err != nil {
				log.Printf("Error creating permission %s: %v", perm.Name, err)
			}
		}
	}

	return nil
}

func GetDB() *gorm.DB {
	db, _ := InitDatabase()
	return db
}

func InitializeRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "Admin", Description: "Administrator with full access"},
		{Name: "Viewer", Description: "User with read-only access"},
		// Add more roles as needed
	}

	for _, role := range roles {
		var existingRole models.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)

		if result.Error != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %v", role.Name, err)
			}
		}
	}

	return nil
}
