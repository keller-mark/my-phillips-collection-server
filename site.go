package main

import (
  "fmt"
  "strconv"
  "log"
  "net/http"
  _ "encoding/json"
  "html/template"
  "github.com/julienschmidt/httprouter"
)


var templates = template.Must(template.ParseFiles(
  "templates/index.html", 
  "templates/survey.html", 
  "templates/header.html", 
  "templates/footer.html",
  "templates/sources.html",
))
var settings = &Settings{Site: "\"Tinder for Museums\""}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  renderTemplate(w, "index", settings)
}

func work(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  settings.WorkID, _ = strconv.Atoi(params.ByName("work_id"))
  /* echo json */
}


func survey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  renderTemplate(w, "survey", settings)
}

func surveySubmit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  /* echo json results of submission */

  fmt.Fprintf(w, "Survey successfully submitted\n")
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
  router.GET("/survey", survey)
  router.POST("/survey", surveySubmit)
  router.GET("/work/:work_id", work)
  router.ServeFiles("/static/*filepath", http.Dir("./static/"))

  log.Fatal(http.ListenAndServe(":80", router))
}
