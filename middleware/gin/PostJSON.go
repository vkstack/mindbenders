package gin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//PostJSONValidator validates a valid JSON in gin.Context
func PostJSONValidator(c *gin.Context) {
	var byteData []byte
	if c.Request.Body != nil {
		byteData, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteData))
	var jsonData interface{}
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
	c.Set("jsonByte", byteData)
}
