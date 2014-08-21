package games

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "server/utils"
  "time"
  // "github.com/go-martini/martini"
)

type Game struct {
  Season  int         `json:"season"`
  Week    int         `json:"week"`
  Date    time.Time   `json:"date"`
     // Week    int   `json:"week"`
     //   Week    int   `json:"week"` 
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  params := r.URL.Query()
  week := params.Get(":week")

  q := datastore.NewQuery("Game").Filter("Week =", week).Limit(32)

  games := make([]Game, 0, 32)
  if _, err := q.GetAll(c, &games); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &games)
}

func AddGameHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var g Game
  if err := utils.ReadJson(r, &g); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  key := datastore.NewIncompleteKey(c, "Game", nil)
  _, err := datastore.Put(c, key, &g)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &g)
}
