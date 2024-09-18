package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Redirect(c *gin.Context) {
	shortURL := c.Param("shortURL")
	URLPair, ok := URLMap[shortURL]

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет урла"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, URLPair.URL.String())
}
