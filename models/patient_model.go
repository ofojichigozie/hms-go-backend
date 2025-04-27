package models

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	RegistrationNumber string    `json:"registrationNumber" gorm:"unique;not null"`
	FirstName          string    `json:"firstName" gorm:"not null"`
	LastName           string    `json:"lastName" gorm:"not null"`
	DateOfBirth        time.Time `json:"dateOfBirth"`
	Gender             string    `json:"gender" gorm:"type:gender_enum;not null"`
	PhoneNumber        string    `json:"phoneNumber" gorm:"not null"`
	Email              string    `json:"email,omitempty"`
	Address            string    `json:"address,omitempty"`
	BloodGroup         string    `json:"bloodGroup,omitempty" gorm:"type:blood_group_enum;default:'unknown'"`
	Genotype           string    `json:"genotype,omitempty" gorm:"type:genotype_enum;default:'unknown'"`
	CreatedBy          uint      `json:"createdBy"`
	UpdatedBy          uint      `json:"updatedBy"`
}

type CreatePatientInput struct {
	FirstName   string `json:"firstName" binding:"required,min=2,max=50"`
	LastName    string `json:"lastName" binding:"required,min=2,max=50"`
	DateOfBirth string `json:"dateOfBirth" binding:"required,datetime=2006-01-02"`
	Gender      string `json:"gender" binding:"required,oneof=male female other"`
	PhoneNumber string `json:"phoneNumber" binding:"required,e164"`
	Email       string `json:"email" binding:"omitempty,email,max=100"`
	Address     string `json:"address" binding:"omitempty,max=200"`
}

type UpdatePatientInput struct {
	FirstName   *string `json:"firstName,omitempty" binding:"omitempty,min=2,max=50"`
	LastName    *string `json:"lastName,omitempty" binding:"omitempty,min=2,max=50"`
	DateOfBirth *string `json:"dateOfBirth,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Gender      *string `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	PhoneNumber *string `json:"phoneNumber,omitempty" binding:"omitempty,e164"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address     *string `json:"address,omitempty" binding:"omitempty,max=200"`
	BloodGroup  *string `json:"bloodGroup,omitempty" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O- unknown"`
	Genotype    *string `json:"genotype,omitempty" binding:"omitempty,oneof=AA AS AC SS SC CC unknown"`
}
