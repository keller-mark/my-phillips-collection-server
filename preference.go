package main

import (
  "github.com/jinzhu/gorm"
)

type Preference struct {
  gorm.Model
  Work_id	  int	`gorm:"column:work_id"`
  User_id	  int	`gorm:"column:user_id"`
  Liked		  int	`gorm:"column:liked"`
}
