package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect set up the database connection
func Connect(host string, user string, password string, dbname string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, dbname)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Couldn't connect to the database")
	}

	return connection
}
