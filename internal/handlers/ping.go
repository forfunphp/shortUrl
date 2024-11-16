package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
)

var db *sql.DB
var logger *zap.Logger

func Ping(c *gin.Context) {

	c.String(http.StatusOK, "OK")
}
