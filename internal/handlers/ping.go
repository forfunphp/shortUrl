package handlers

import (
	"database/sql"
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var db *sql.DB
var err error

func Ping(c *gin.Context) {
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
	c.String(http.StatusOK, "OK")
}
