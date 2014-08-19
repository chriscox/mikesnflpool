package games

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "server/utils"
)

type Game struct {
  Season int
  Week int
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Game").Limit(32)
  games := make([]Game, 0, 32)
  if _, err := q.GetAll(c, &games); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &games)
}

func AddGameHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  g := Game {
          Season: 2014,
          Week: 1,
        }

  key := datastore.NewIncompleteKey(c, "Game", nil)
  _, err := datastore.Put(c, key, &g)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}
