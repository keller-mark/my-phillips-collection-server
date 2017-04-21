package main

import (
  "github.com/jinzhu/gorm"
)

type Work struct {
  gorm.Model
  Name	  string	`gorm:"column:name"`
  Artist  string	`gorm:"column:artist"`
}
