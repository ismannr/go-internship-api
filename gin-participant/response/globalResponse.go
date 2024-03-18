package response

import (
	"github.com/gin-gonic/gin"
	"time"
)

func GlobalResponse(c *gin.Context, message string, status int, data interface{}) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	timestamp := time.Now().In(loc).Format("2006-01-02T15:04:05.000Z07:00")

	response := gin.H{
		"status":    status,
		"message":   message,
		"timestamp": timestamp,
		"data":      data,
	}

	// Set response status code and JSON body
	c.JSON(status, response)
}
