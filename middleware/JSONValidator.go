package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidatePostJSON(c *gin.Context) {
	var jsonData interface{}
	byteData, _ := c.GetRawData()
	err := json.Unmarshal(byteData, &jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"error":   error.Error(err),
			"message": "Bad JSON.",
		})
		c.Abort()
		return
	}
	c.Set("json", jsonData)
	c.Set("jsonByte", byteData)
	c.Next()
}
