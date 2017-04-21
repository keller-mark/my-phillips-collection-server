package main

import ( 
  "os"
  "bytes"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

func DB() *gorm.DB {
  var buffer bytes.Buffer
  buffer.WriteString(os.Getenv("MYSQL_PHILLIPS_USERNAME"))
  buffer.WriteString(":")
  buffer.WriteString(os.Getenv("MYSQL_PHILLIPS_PASSWORD"))
  buffer.WriteString("@tcp(localhost:3306)/phillips_data?parseTime=true")
  
  db, err := gorm.Open("mysql", buffer.String())
  if err != nil {
    panic("failed to connect to database")
  }
  return db
}
