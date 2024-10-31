package handlers

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	gzip *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzip.Write(b)
}

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		originalResponseWriter := c.Writer
		c.Writer = &gzipResponseWriter{Writer: c.Writer, gzip: gzip.NewWriter(c.Writer)}

		c.Writer.Header().Set("Content-Encoding", "gzip")

		c.Next()

		if gw, ok := c.Writer.(*gzipResponseWriter); ok {
			gw.gzip.Close()
		}

		c.Writer = originalResponseWriter
	}
}
