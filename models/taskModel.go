package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Description string
	IsComplete  bool
}
