package main

import (
"fmt"
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
  "templates/visualize-likes.html",
  "templates/header.html", 
  "templates/footer.html",
  "templates/sources.html",
))
var settings = &Settings{Site: "My Phillips Collection"}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  renderTemplate(w, "index", settings)
}

func work(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  db := DB()
  scannedWork, _ := strconv.Atoi(params.ByName("work_id"))
  foundWork := []Work{}
  db.Table("works").Find(&foundWork, "id = ?", scannedWork)
  //  fmt.Println(foundWork[0].Title)
  /* echo json */
  w.Header().Set("Content-Type", "application/json")
  jData, err := json.Marshal(foundWork[0])
  if err != nil {
    panic(err)
    return
  }
  w.Write(jData)
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

func newPreference(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  db := DB()
  user_id, _ := strconv.Atoi(r.FormValue("user_id"))
  work_id, _ := strconv.Atoi(r.FormValue("work_id"))
  liked, _ := strconv.Atoi(r.FormValue("liked"))
//  fmt.Println(user_id, work_id, liked)

  preferenceSubmit := Preference{
    User_id: user_id,
    Work_id: work_id,
    Liked: liked,
  }
//  fmt.Println(preferenceSubmit)
  foundPreferences := []Preference{}
  db.Table("preferences").Find(&foundPreferences, "user_id = ? AND work_id = ?", user_id, work_id)
//  fmt.Println(foundPreferences)
  if(len(foundPreferences) == 0){
    db.Create(&preferenceSubmit)
  } else if (foundPreferences[0].Liked != liked) {
    foundPreferences[0].Liked = liked
    db.Save(&foundPreferences[0])
  }
}

func visualizeLikes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  db := DB()
  
  type WorkData struct {
    TheWork     Work
    Likes	int
    Dislikes	int
  }
  allWorkData := make(map[string]WorkData)

  preferences := []Preference{}
  db.Model(&Preference{}).Find(&preferences);
  for _, preference := range preferences {
    if _, ok := allWorkData[fmt.Sprint(preference.Work_id)]; !ok { 
      work := Work{}
      db.Model(&Work{}).Where("id = ?", preference.Work_id).Find(&work)
      workData := WorkData{}
      workData.TheWork = work;
      if preference.Liked == 1 {
	workData.Likes = 1
	workData.Dislikes = 0
      } else {
	workData.Likes = 0
	workData.Dislikes = 1
      }
      allWorkData[fmt.Sprint(work.ID)] = workData
    } else {
      
      workData := allWorkData[fmt.Sprint(preference.Work_id)]
      if preference.Liked == 1 {
	workData.Likes++
      } else {
	workData.Dislikes++
      }
      allWorkData[fmt.Sprint(preference.Work_id)] = workData
    }
  }

  b, _ := json.Marshal(allWorkData)
  
  w.Header().Set("Content-Type", "application/json")
  w.Write(b)
}

func main() { 
  router := httprouter.New()

  router.GET("/", home)
  router.GET("/survey", survey)
  router.POST("/survey", surveySubmit)
  router.GET("/work/:work_id", work)
  router.POST("/work/preference", newPreference)
  router.GET("/visualize/likes", visualizeLikes)
  router.ServeFiles("/static/*filepath", http.Dir("./static/"))

  // ParseWorksCSV("../../../perm_coll_filtered_20170201.csv")

  log.Fatal(http.ListenAndServe(":80", router))
}
