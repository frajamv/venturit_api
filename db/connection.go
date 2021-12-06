package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")
var ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
	os.Getenv("db_user"),
	os.Getenv("db_pass"),
	os.Getenv("db_host"),
	os.Getenv("db_port"),
	os.Getenv("db_name"))

func GetDBconnection() (*sql.DB, error) {
	return sql.Open("mysql", ConnectionString)
}
