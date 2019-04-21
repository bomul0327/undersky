package undersky

import (
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	newDB, err := gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port=5432 user="+os.Getenv("DB_USERNAME")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
	if err != nil {
		panic(err)
	}

	maxIdleConn, err := strconv.Atoi(GetenvOrDefault("DATABASE_MAX_IDLE_CONN", "20"))
	if err != nil {
		panic(err)
	}

	newDB.DB().SetMaxIdleConns(maxIdleConn)
	DB = newDB
}
