package responses

import "github.com/gin-gonic/gin"

type StandardResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Status:  false,
		Message: message,
		Data:    data,
	})
}
