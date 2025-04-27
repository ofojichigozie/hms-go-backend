package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetCurrentStaff(c *gin.Context) (CurrentStaff, error) {
	val, exists := c.Get(CurrentStaffKey)
	if !exists {
		return CurrentStaff{}, errors.New("staff context missing")
	}

	user, ok := val.(CurrentStaff)
	if !ok {
		return CurrentStaff{}, errors.New("invalid staff context type")
	}

	if user.ID == 0 {
		return CurrentStaff{}, errors.New("invalid staff ID")
	}

	return user, nil
}
