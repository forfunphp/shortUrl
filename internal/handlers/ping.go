package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"net/http"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database") // Замените на ваше подключение
	if err != nil {
		panic(err)
	}
	// Проверка подключения (необязательно, но рекомендуется)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func Ping(c *gin.Context) {
	err := db.Ping()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "OK")
}
