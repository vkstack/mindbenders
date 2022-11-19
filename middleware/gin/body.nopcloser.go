package gin

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func BodyReaderNop(c *gin.Context) {
	if c.Request.Body != nil {
		buff := new(bytes.Buffer)
		buff.ReadFrom(c.Request.Body)
		c.Request.Body = io.NopCloser(buff)
	}
}
