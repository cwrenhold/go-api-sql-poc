package main

import (
	"github.com/cwrenhold/go-api-sql-poc/initializers"
	"github.com/cwrenhold/go-api-sql-poc/models"
)

func init() {
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Task{})
}
