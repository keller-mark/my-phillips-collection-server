package main

import (
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  survey_id   int `gorm:"column:survey_id"`
}
