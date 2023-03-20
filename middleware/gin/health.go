package gin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func GetHealthHandlerWithDB(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats := gin.H{
			"go-routines": runtime.NumGoroutine(),
			"MySQL":       db.Stats(),
		}
		x, _ := json.Marshal(stats)
		c.JSON(http.StatusOK, stats)
	}
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
