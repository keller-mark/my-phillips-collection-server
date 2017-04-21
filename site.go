package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "html/template"
  "github.com/julienschmidt/httprouter"
)


var templates = template.Must(template.ParseFiles("templates/index.html"))
var settings = &Settings{Site: "Tinder for Museums"}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  renderTemplate(w, "index", settings)
}

func renderTemplate(w http.ResponseWriter, tmpl string, settings *Settings) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", settings)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}


func main() { 
  router := httprouter.New()
  router.GET("/", home)
  router.ServeFiles("/static/*filepath", http.Dir("./static/"))

  db := DB()
  defer db.Close()

  var user User
  db.First(&user)
  jsonUser1, _ := json.Marshal(user)
  fmt.Println(string(jsonUser1))


  log.Fatal(http.ListenAndServe(":80", router))
}
