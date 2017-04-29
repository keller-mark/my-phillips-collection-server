package main

import (
  _ "fmt"
  "strconv"
  "log"
  "net/http"
  "encoding/json"
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
var settings = &Settings{Site: "Project Name"}

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
  db := DB()
  age, _ := strconv.Atoi(r.FormValue("visitor_age"))
  gender, _ := strconv.Atoi(r.FormValue("visitor_gender"))
  country := r.FormValue("visitor_country")
  state := r.FormValue("visitor_state")

  user := User{
    Age: age, 
    Gender: gender, 
    Country: country, 
    State: state,
  }
  db.Create(&user)
  
  var surveyResult struct {
    Success bool
    Message string
    UserID  uint
  }
  surveyResult.Success = true
  surveyResult.Message = "Survey successfully submitted"
  surveyResult.UserID = user.ID
  
  b, err := json.Marshal(surveyResult)
  if err != nil {
    log.Fatal(err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(b)
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

  // ParseWorksCSV("../../../perm_coll_filtered_20170201.csv")

  log.Fatal(http.ListenAndServe(":80", router))
}
