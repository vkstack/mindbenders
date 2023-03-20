package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) { c.JSON(http.StatusOK, "OK") }
