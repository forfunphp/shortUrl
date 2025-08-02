package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Redirect(c *gin.Context) {
	shortURL := c.Param("shortURL")
	URLPair, ok := URLMap[shortURL]

	filePath := Cfg.EnvFilePath

	logger2, _ := zap.NewDevelopment()
	defer logger2.Sync()
	logger2.Info("Request proceju00003",
		zap.String("fullURL", shortURL),  // Добавляем полный URL
		zap.String("filePath", filePath), // Добавляем полный URL
	)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет урла"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, URLPair.URL.String())
}
