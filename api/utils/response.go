package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendSuccessJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func SendMessageWithStatus(c *gin.Context, message string, status int) {
	c.JSON(status, gin.H{
		"code":    status,
		"message": message,
	})
}
