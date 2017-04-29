package main

import (
  "encoding/csv"
  "os"
  "io"
  "log"
  "unicode"
  "strings"
  "bytes"
  "golang.org/x/text/unicode/norm"
  "golang.org/x/text/transform"
  "golang.org/x/text/runes"
  "github.com/jinzhu/gorm"
)

func ParseWorksCSV(file string) {
  f, err := os.Open(file)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  csvReader := csv.NewReader(f)
  
  headMap := make(map[string]int)
  
  db := DB()

  for row := 0; ; row++ {
    record, err := csvReader.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      log.Fatal(err)
    }
    if row == 0 {
      for col := 0; col < len(record); col++ {
	headMap[record[col]] = col
      }
    } else {
      createWork(db, &headMap, record)
    }


  }
}

func createWork(db *gorm.DB, pHeadMap *map[string]int, record []string) { 
    work := Work{}
    
    headMap := *pHeadMap

    for _, v := range headMap {
      record[v] = removeAccents(record[v])
      record[v] = cleanString(record[v])
    }

    work.PhillipsID	  = record[headMap["phillips_id"]]
    work.Title		  = record[headMap["title"]]
    work.Maker		  = record[headMap["maker"]]
    work.DateMade	  = record[headMap["date_made"]]
    work.Culture	  = record[headMap["culture"]]
    work.Materials	  = record[headMap["materials"]]
    work.CreditLine	  = record[headMap["credit_line"]]
    work.ItemName	  = record[headMap["item_name"]]
    work.Movement	  = record[headMap["movement"]]
    work.Century	  = record[headMap["century"]]
    work.Lifespan	  = record[headMap["lifespan"]]
    work.Continent	  = record[headMap["continent"]]
    work.Gender		  = record[headMap["gender"]]
    work.Year		  = record[headMap["year"]]
    
    db.Where(Work{PhillipsID: work.PhillipsID}).FirstOrCreate(&work)
}

func removeAccents(s string) string {
  t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
  r, _, _ := transform.String(t, s)
  
  return string(r)
}

func cleanString(s string) string {
  runes := bytes.Runes([]byte(s))
  var result string = s
  if len(runes) > 255 {
    result = string(runes[:255])
  }
  return strings.Replace(result, "\xEF\xBF\xBD", "", -1)
}
