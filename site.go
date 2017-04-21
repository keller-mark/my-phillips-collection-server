package main

import (
  _ "fmt"
  "strconv"
  "log"
  "net/http"
  _ "encoding/json"
  "html/template"
  "github.com/julienschmidt/httprouter"
)


var templates = template.Must(template.ParseFiles("templates/index.html", "templates/survey.html"))
var settings = &Settings{Site: "Tinder for Museums"}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  renderTemplate(w, "index", settings)
}

func survey(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  settings.UserID, _ = strconv.Atoi(params.ByName("user_id"))
  renderTemplate(w, "survey", settings)
}

func surveySubmit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  renderTemplate(w, "survey", settings)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", data)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}


func main() { 
  router := httprouter.New()

  router.GET("/", home)
  router.GET("/survey/:user_id", survey)
  router.POST("/survey/:user_id", surveySubmit)

  router.ServeFiles("/static/*filepath", http.Dir("./static/"))

  log.Fatal(http.ListenAndServe(":80", router))
}
