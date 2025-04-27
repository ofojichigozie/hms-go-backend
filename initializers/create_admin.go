package initializers

import (
	"errors"
	"fmt"
	"log"

	"github.com/ofojichigozie/hms-go-backend/constants"
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/utils"
	"gorm.io/gorm"
)

const (
	AdminEmail      = "system.admin@hospital.com"
	AdminPassword   = "Admin@12345"
	AdminEmployeeID = "ADM0001"
	AdminFirstName  = "System"
	AdminLastName   = "Admin"
)

func CreateInitialAdmin(db *gorm.DB) (*models.Staff, error) {
	var existing models.Staff
	err := db.Where("email = ?", AdminEmail).First(&existing).Error

	if err == nil {
		log.Printf("Admin account already exists (ID: %d, Email: %s)", existing.ID, existing.Email)
		return &existing, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashedPassword, err := utils.HashPassword(AdminPassword)
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		admin := models.Staff{
			EmployeeID:   AdminEmployeeID,
			FirstName:    AdminFirstName,
			LastName:     AdminLastName,
			Email:        AdminEmail,
			PasswordHash: hashedPassword,
			Role:         constants.Roles.ADMIN,
			IsActive:     true,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to create admin in database: %v", err)
			return nil, fmt.Errorf("failed to create admin: %w", err)
		}

		log.Printf("Successfully created initial admin (ID: %d, Email: %s)", admin.ID, admin.Email)
		return &admin, nil
	}

	log.Printf("Database error while checking for admin: %v", err)
	return nil, fmt.Errorf("database error checking for admin: %w", err)
}
