package utils

import (
	"fmt"
	"time"
)

func GeneratePatientID() string {
	return fmt.Sprintf("PAT-%d-%04d", time.Now().Year(), time.Now().UnixNano()%10000)
}
