package main

import (
	"github.com/cwrenhold/go-api-sql-poc/initializers"
	"github.com/cwrenhold/go-api-sql-poc/models"
)

func init() {
	initializers.GormConnectToDB()
}

func main() {
	initializers.GormDB.AutoMigrate(&models.Task{})

	// If there aren't any tasks in the database, add some
	var count int64
	initializers.GormDB.Model(&models.Task{}).Count(&count)

	if count == 0 {
		defaultTasks := []models.Task{
			{Description: "Buy groceries", IsComplete: false},
			{Description: "Do laundry", IsComplete: true},
		}

		for _, task := range defaultTasks {
			initializers.GormDB.Create(&task)
		}
	}
}
