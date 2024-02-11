package db

import (
	"dirwatcher/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

/*
*

	This method will create a connection object to postgres Sql
	@param void
	@return error
*/
func Connect() error {

	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("[db.Connect]Error while opening a connection : %v", err)
		return err
	}

	db.AutoMigrate(&models.Configurations{}, &models.Task{})
	DB = db
	return nil
}
