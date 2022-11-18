package gin

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func BodyReaderNop(c *gin.Context) {
	if c.Request.Body != nil {
		raw, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
	}
}
