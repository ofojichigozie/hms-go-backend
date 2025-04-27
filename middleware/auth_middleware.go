package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ofojichigozie/hms-go-backend/responses"
	"github.com/ofojichigozie/hms-go-backend/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("accessToken")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			responses.Error(ctx, http.StatusUnauthorized, "Unauthorized", "An access token is required")
			ctx.Abort()
			return
		}

		user, err := utils.VerifyToken(accessToken)
		if err != nil {
			responses.Error(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
			ctx.Abort()
			return
		}

		ctx.Set(CurrentStaffKey, FromJWTClaims(user))
		ctx.Next()
	}
}

func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		staff, err := GetCurrentStaff(ctx)
		if err != nil {
			responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
			ctx.Abort()
			return
		}

		if !slices.Contains(allowedRoles, staff.Role) {
			responses.Error(ctx, http.StatusForbidden,
				"Access denied", "Required role(s): "+strings.Join(allowedRoles, ", "))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
