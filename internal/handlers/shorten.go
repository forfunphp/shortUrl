package handlers

import (
	"github.com/gin-gonic/gin"
)

func Shorten(c *gin.Context) {
	c.Request.URL.Path = "/"
	ReduceURL(c) // Используем существующий ReduceURL
}
