package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	PatientID      uint      `json:"patientId" gorm:"not null"`
	ReceptionistID uint      `json:"receptionistId" gorm:"not null"`
	DoctorID       *uint     `json:"doctorId,omitempty"`
	Department     string    `json:"department" gorm:"type:department_type;not null"`
	ScheduledAt    time.Time `json:"scheduledAt" gorm:"not null"`
	Duration       int       `json:"duration" gorm:"default:30"`
	Status         string    `json:"status" gorm:"type:appointment_status;default:'scheduled'"`
	Reason         string    `json:"reason,omitempty" gorm:"size:500"`
	UpdatedBy      uint      `json:"updatedBy"`
}

type CreateAppointmentInput struct {
	PatientID  uint   `json:"patientId" binding:"required"`
	Department string `json:"department" binding:"required,oneof=general cardiology pediatrics"`
	Duration   int    `json:"duration" binding:"omitempty,min=15,max=120"`
	Reason     string `json:"reason" binding:"omitempty,max=500"`
}

type UpdateAppointmentInput struct {
	DoctorID   *uint   `json:"doctorId,omitempty"`
	Department *string `json:"department,omitempty" binding:"omitempty,oneof=general cardiology pediatrics"`
	Status     *string `json:"status,omitempty" binding:"omitempty,oneof=scheduled completed cancelled no_show"`
	Reason     *string `json:"reason,omitempty" binding:"omitempty,max=1000"`
}
