package main

import (
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Age	    int	    `gorm:"column:age"`
  Gender    int	    `gorm:"column:gender"`
  Country   string  `gorm:"column:country"`
  State	    string  `gorm:"column:state"`
}
