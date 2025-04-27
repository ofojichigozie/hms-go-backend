package models

import (
	"time"

	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	EmployeeID     string     `json:"employeeId" gorm:"unique;not null"`
	FirstName      string     `json:"firstName" gorm:"not null"`
	LastName       string     `json:"lastName" gorm:"not null"`
	PhoneNumber    string     `json:"phoneNumber" gorm:"not null"`
	Email          string     `json:"email" gorm:"unique;not null"`
	PasswordHash   string     `json:"-" gorm:"not null"`
	Role           string     `json:"role" gorm:"type:role_enum;not null"`
	IsActive       bool       `json:"isActive" gorm:"default:true"`
	LastLogin      *time.Time `json:"lastLogin,omitempty"`
	LicenseNumber  *string    `json:"licenseNumber,omitempty" gorm:"unique"`
	Specialization *string    `json:"specialization,omitempty"`
	Department     *string    `json:"department,omitempty"`
}

type CreateStaffInput struct {
	EmployeeID     string  `json:"employeeId" binding:"required"`
	FirstName      string  `json:"firstName" binding:"required"`
	LastName       string  `json:"lastName" binding:"required"`
	PhoneNumber    string  `json:"phoneNumber" binding:"required"`
	Email          string  `json:"email" binding:"required,email"`
	Password       string  `json:"password" binding:"required,min=8"`
	Role           string  `json:"role" binding:"required,oneof=admin doctor receptionist"`
	LicenseNumber  *string `json:"licenseNumber,omitempty" binding:"required_if=Role doctor"`
	Specialization *string `json:"specialization,omitempty"`
	Department     *string `json:"department,omitempty"`
}

type UpdateStaffInput struct {
	FirstName      *string `json:"firstName,omitempty" binding:"omitempty,min=2"`
	LastName       *string `json:"lastName,omitempty" binding:"omitempty,min=2"`
	PhoneNumber    *string `json:"phoneNumber,omitempty"`
	Email          *string `json:"email,omitempty" binding:"omitempty,email"`
	IsActive       *bool   `json:"isActive,omitempty"`
	LicenseNumber  *string `json:"licenseNumber,omitempty" binding:"omitempty,required_if=Role doctor"`
	Specialization *string `json:"specialization,omitempty"`
	Department     *string `json:"department,omitempty"`
}
