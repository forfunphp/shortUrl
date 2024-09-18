package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Redirect(c *gin.Context) {
	shortURL := c.Param("shortURL")
	urlPair, ok := urlMap[shortURL]

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет урла"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, urlPair.URL.String())
}
