package main

import (
  "github.com/jinzhu/gorm"
)

type Work struct {
  gorm.Model
  PhillipsID	  string	`gorm:"column:phillips_id"`
  Title		  string	`gorm:"column:title"`
  Maker		  string	`gorm:"column:maker"`
  DateMade	  string	`gorm:"column:date_made"`
  Culture	  string	`gorm:"column:culture"`
  Materials	  string	`gorm:"column:materials"`
  CreditLine	  string	`gorm:"column:credit_line"`
  ItemName	  string	`gorm:"column:item_name"`
  Movement	  string	`gorm:"column:movement"`
  Century	  string	`gorm:"column:century"`
  Lifespan	  string	`gorm:"column:lifespan"`
  Continent	  string	`gorm:"column:continent"`
  Gender	  string	`gorm:"column:gender"`
  Year		  string	`gorm:"column:year"`
}
