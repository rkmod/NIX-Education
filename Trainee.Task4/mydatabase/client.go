package mydatabase

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Connector *gorm.DB

// NewConnection opens connection to DB and returns it
func NewDBConnection(connectionString string) error {
	var err error

	Connector, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("Successfully connected to the database.")
	return nil
}
