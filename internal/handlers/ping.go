package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Импорт драйвера MySQL
	"go.uber.org/zap"
)

var db *sql.DB
var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize logger: %v", err)
	}
	defer logger.Sync()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsnPtr := flag.String("d", "", "MySQL DSN (database source name)")
		flag.Parse()
		dsn = *dsnPtr
	}

	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable or -d flag is required.")
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Ошибка при подключении к базе данных", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("Ошибка при проверке подключения к базе данных", zap.Error(err))
	}
	logger.Info("Подключение к базе данных установлено.")
}

func Ping(c *gin.Context) {
	err := db.Ping()
	if err != nil {
		logger.Error("Ошибка при проверке подключения к базе данных", zap.Error(err))
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "OK")
}
