package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func HandleError(c *gin.Context, status int, err error, details string) {
	c.JSON(status, ErrorResponse{
		Error:   err.Error(),
		Details: details,
	})
}
