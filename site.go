package main

import (
  "fmt"
  "bytes"
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

func visualizeWorkList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  db := DB();

  type WorkListing struct {
    ID	  string
    Title string
  }
  allWorks := []Work{}
  db.Find(&allWorks)
  
  workListings := []WorkListing{}

  for _, work := range allWorks {
    workListings = append(workListings, WorkListing{fmt.Sprint(work.ID), work.Title})
  }
  b, _ := json.Marshal(workListings)

  w.Header().Set("Content-Type", "application/json")
  w.Write(b)
}

func visualizeWork(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  db := DB();
  
  work_id := params.ByName("work_id")
  preferences := []Preference{}
  db.Model(&Preference{}).Where("work_id = ?", work_id).Find(&preferences)

  type WorkStats struct {
    Likes		int
    Dislikes		int
    Net			int
    
    LikesByAge		map[string]int
    DislikesByAge	map[string]int
    NetByAge		map[string]int
    
    LikesByGender	map[string]int
    DislikesByGender	map[string]int
    NetByGender		map[string]int
    
    LikesByLocation	map[string]int
    DislikesByLocation	map[string]int
    NetByLocation	map[string]int
  }

  workStats := WorkStats{}
  workStats.LikesByAge = make(map[string]int)
  workStats.DislikesByAge = make(map[string]int)
  workStats.NetByAge = make(map[string]int)

  workStats.LikesByGender = make(map[string]int)
  workStats.DislikesByGender = make(map[string]int)
  workStats.NetByGender = make(map[string]int)

  workStats.LikesByLocation = make(map[string]int)
  workStats.DislikesByLocation = make(map[string]int)
  workStats.NetByLocation = make(map[string]int)
  
  for _, pref := range preferences {
    user := User{}
    db.First(&user, pref.User_id)
    var location bytes.Buffer
    location.WriteString(user.Country)
    if user.Country == "US" {
      location.WriteString("-")
      location.WriteString(user.State)
    }
    locationString := location.String()

    if pref.Liked == 1 {
      workStats.Likes++
      if _, ok := workStats.LikesByAge[fmt.Sprint(user.Age)]; ok {
	workStats.LikesByAge[fmt.Sprint(user.Age)]++
      } else {
	workStats.LikesByAge[fmt.Sprint(user.Age)] = 1
      }

      if _, ok := workStats.LikesByGender[fmt.Sprint(user.Gender)]; ok {
	workStats.LikesByGender[fmt.Sprint(user.Gender)]++
      } else {
	workStats.LikesByGender[fmt.Sprint(user.Gender)] = 1
      }
      
      if _, ok := workStats.LikesByLocation[locationString]; ok {
	workStats.LikesByLocation[locationString]++
      } else {
	workStats.LikesByLocation[locationString] = 1
      }
    } else {
      workStats.Dislikes++
      if _, ok := workStats.DislikesByAge[fmt.Sprint(user.Age)]; ok {
	workStats.DislikesByAge[fmt.Sprint(user.Age)]++
      } else {
	workStats.DislikesByAge[fmt.Sprint(user.Age)] = 1
      }

      if _, ok := workStats.DislikesByGender[fmt.Sprint(user.Gender)]; ok {
	workStats.DislikesByGender[fmt.Sprint(user.Gender)]++
      } else {
	workStats.DislikesByGender[fmt.Sprint(user.Gender)] = 1
      }
      
      if _, ok := workStats.DislikesByLocation[locationString]; ok {
	workStats.DislikesByLocation[locationString]++
      } else {
	workStats.DislikesByLocation[locationString] = 1
      }
    }
  }
  workStats.Net = workStats.Likes - workStats.Dislikes
  for i := 1; i <= 5; i++ {
    num_likes := 0
    if _, ok := workStats.LikesByAge[fmt.Sprint(i)]; ok {
      num_likes = workStats.LikesByAge[fmt.Sprint(i)]
    }
    num_dislikes := 0
    if _, ok := workStats.DislikesByAge[fmt.Sprint(i)]; ok {
      num_dislikes = workStats.DislikesByAge[fmt.Sprint(i)]
    }
    workStats.NetByAge[fmt.Sprint(i)] = num_likes - num_dislikes
  }
  for i := 1; i <= 3; i++ {
    num_likes := 0
    if _, ok := workStats.LikesByGender[fmt.Sprint(i)]; ok {
      num_likes = workStats.LikesByGender[fmt.Sprint(i)]
    }
    num_dislikes := 0
    if _, ok := workStats.DislikesByGender[fmt.Sprint(i)]; ok {
      num_dislikes = workStats.DislikesByGender[fmt.Sprint(i)]
    }
    workStats.NetByGender[fmt.Sprint(i)] = num_likes - num_dislikes
  }
  for key, val := range workStats.LikesByLocation {
    workStats.NetByLocation[key] = val
  }
  for key, val := range workStats.DislikesByLocation {
    if _, ok := workStats.NetByLocation[key]; ok {
      workStats.NetByLocation[key] -= val
    } else {
      workStats.NetByLocation[key] = -(val)
    }
  }
  
  b, err := json.Marshal(workStats)
  if err != nil {
    log.Fatal(err)
  } 
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
  router.GET("/visualize/work-list", visualizeWorkList)
  router.GET("/visualize/work/:work_id", visualizeWork)

  router.ServeFiles("/static/*filepath", http.Dir("./static/"))

  // ParseWorksCSV("../../../perm_coll_filtered_20170201.csv")

  log.Fatal(http.ListenAndServe(":80", router))
}
