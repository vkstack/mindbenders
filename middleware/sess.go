package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func SessionManager(c *gin.Context) {
	if sess, err := c.Cookie("sessid"); err != nil || len(sess) == 0 {
		c.SetCookie("sessid", xid.New().String(), int(4*time.Hour), "/", c.Request.URL.Host, true, false)
	}
	c.Next()
}
