package database

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

type Config struct {
	URL       string
	MaxDBConn int
}

func Init(config *Config) error {
	var err error
	once.Do(func() {
		var sqlDB *sql.DB
		sqlDB, err = sql.Open("postgres", config.URL)
		if err != nil {
			log.Println("Unable to open postges connection. Err:", err)
			return
		}

		sqlDB.SetMaxIdleConns(config.MaxDBConn)
		sqlDB.SetMaxOpenConns(config.MaxDBConn)
		sqlDB.SetConnMaxLifetime(time.Hour)

		DB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Println("Unable to open postges gorm connection. Err:", err)
			return
		}

		log.Println("Successfully established database connection")
	})

	return err
}

type DBConn struct {
	*gorm.DB
}

func New() *DBConn {
	return &DBConn{
		DB: DB,
	}
}
