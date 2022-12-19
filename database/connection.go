package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(host string, user string, dbname string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=docker dbname=%s", host, user, dbname)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Couldn't connect to the database")
	}

	return connection
}
