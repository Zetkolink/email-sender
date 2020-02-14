package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var connection *gorm.DB

func connectToDb(dbUrl string, schema string) *gorm.DB {
	db, err := gorm.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	gorm.DefaultTableNameHandler = func(dbVeiculosGorm *gorm.DB, defaultTableName string) string {
		return schema + "." + defaultTableName
	}

	return db
}

func ConnectToDb(dbUrl string, schema string) (*gorm.DB, error) {
	if connection == nil {
		connection = connectToDb(dbUrl, schema)
	}
	check := connection.DB().Ping()
	for check != nil {
		log.Println("Connection has been lost. Attempt to reconnect after 5 seconds...")
		_ = connection.Close()
		time.Sleep(5 * time.Second)
		log.Println("Reconnecting...")
		connection = connectToDb(dbUrl, schema)
		check = connection.DB().Ping()
	}
}
