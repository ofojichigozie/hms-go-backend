package middleware

import "github.com/ofojichigozie/hms-go-backend/utils"

type CurrentStaff struct {
	ID   uint   `json:"userId"`
	Role string `json:"role"`
}

const (
	CurrentStaffKey = "currentStaff"
)

func FromJWTClaims(claims *utils.JWTClaims) CurrentStaff {
	return CurrentStaff{
		ID:   claims.ID,
		Role: claims.Role,
	}
}
