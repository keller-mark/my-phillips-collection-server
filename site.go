package main

import (
  _ "fmt"
  "log"
  "net/http"
  _ "encoding/json"
  "html/template"
  "github.com/julienschmidt/httprouter"
  _ "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Page struct {
  Title string
}

var templates = template.Must(template.ParseFiles("templates/index.html"))

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  p := &Page{Title: "Home"}
  renderTemplate(w, "index", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func main() {
  router := httprouter.New()
  router.GET("/", home)
  router.ServeFiles("/static/*filepath", http.Dir("./static/"))

  log.Fatal(http.ListenAndServe(":80", router))
}
