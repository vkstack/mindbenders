package gin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func SessionManager(c *gin.Context) {
	if _, err := c.Cookie("sessid"); err != nil {
		c.SetCookie("sessid", xid.New().String(), int(4*time.Hour), "/", c.Request.URL.Host, true, false)
	}
	c.Next()
}
