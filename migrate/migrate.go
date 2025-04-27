package main

import (
	"fmt"

	"github.com/ofojichigozie/hms-go-backend/initializers"
	"github.com/ofojichigozie/hms-go-backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	createEnums()

	err := initializers.DB.AutoMigrate(&models.Staff{}, &models.Patient{},
		&models.Appointment{}, &models.ClinicalNote{})
	if err != nil {
		panic("Migration failed: " + err.Error())
	}
	fmt.Println("Database migration successful")
}

func createEnums() {
	DB := initializers.DB

	DB.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_enum') THEN
			CREATE TYPE role_enum AS ENUM ('admin', 'doctor', 'receptionist');
		END IF;
	END
	$$;`)

	DB.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender_enum') THEN
			CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');
		END IF;
	END
	$$;`)

	DB.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blood_group_enum') THEN
			CREATE TYPE blood_group_enum AS ENUM ('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-', 'unknown');
		END IF;
	END
	$$;`)

	DB.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'genotype_enum') THEN
			CREATE TYPE genotype_enum AS ENUM ('AA', 'AS', 'AC', 'SS', 'SC', 'CC', 'unknown');
		END IF;
	END
	$$;`)

	DB.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'department_type') THEN
			CREATE TYPE department_type AS ENUM (
				'general', 
				'cardiology', 
				'pediatrics', 
				'orthopedics', 
				'neurology', 
				'dermatology', 
				'psychiatry', 
				'oncology', 
				'gynecology', 
				'endocrinology'
			);
		END IF;
	END
	$$;
	`)

	DB.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'appointment_status') THEN
			CREATE TYPE appointment_status AS ENUM (
				'scheduled', 
				'completed', 
				'cancelled', 
				'no_show'
			);
		END IF;
	END
	$$;
	`)
}
