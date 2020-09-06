package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func SessionManager(c *gin.Context) {
	if sess, err := c.Cookie("_sessid"); err != nil || len(sess) == 0 {
		c.SetCookie("_sessid", xid.New().String(), 2, "/", c.Request.URL.Host, true, false)
	}
	c.Next()
}
