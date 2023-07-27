package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cwrenhold/go-api-sql-poc/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func GormConnectToDB() {
	host := os.Getenv("POSTGRES_HOSTNAME")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London", host, user, password, dbname, port)
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}
}

func init() {
	GormConnectToDB()
}

func main() {
	GormDB.AutoMigrate(&models.Task{})

	// If there aren't any tasks in the database, add some
	var count int64
	GormDB.Model(&models.Task{}).Count(&count)

	if count == 0 {
		defaultTasks := []models.Task{
			{Description: "Buy groceries", IsComplete: false},
			{Description: "Do laundry", IsComplete: true},
		}

		for _, task := range defaultTasks {
			GormDB.Create(&task)
		}
	}
}
