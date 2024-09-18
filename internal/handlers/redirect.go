package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Redirect(c *gin.Context) {
	shortURL := c.Param("shortURL")
	UrlPair, ok := UrlMap[shortURL]

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет урла"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, UrlPair.URL.String())
}
