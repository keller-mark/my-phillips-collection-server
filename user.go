package main

import (
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Age	    int	    `gorm:"column:age"`
  Gender    int	    `gorm:"column:gender"`
  Location  string  `gorm:"column:location"`
}
