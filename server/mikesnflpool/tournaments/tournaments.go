package tournaments

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "github.com/go-martini/martini"
  "server/mikesnflpool/utils"
  "strconv"
)

type Tournament struct {
  Name          string          `json:"name"`
  Season        int             `json:"season"`
  TournamentKey *datastore.Key  `json:"tournamentKey" datastore:"-"`
}

type TournamentUser struct {
  UserKey       *datastore.Key  `json:"userKey"`
  Admin         bool            `json:"-"`
}

type GameEvent struct {
  Season        int             `json:"season"`
  Week          int             `json:"week"`
}

func TournamentHandler(parms martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  season,_ := strconv.Atoi(parms["season"])

  // Get tournaments
  q := datastore.NewQuery("Tournament").Filter("Season =", season)
  var tournaments []Tournament
  keys, err := q.GetAll(c, &tournaments); if err != nil {
    panic(err.Error)
  }

  // Associate keys with tournament
  for i := range keys {
    t := &tournaments[i]
    t.TournamentKey = keys[i]
  }
  utils.ServeJson(w, &tournaments)
}

func AddTournamentHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t Tournament
  if err := utils.ReadJson(r, &t); err != nil {
    panic(err.Error())
  }

  // Add tournament
  key := datastore.NewIncompleteKey(c, "Tournament", nil)
  tournyKey, err := datastore.Put(c, key, &t); if err != nil {
    panic(err.Error)
  }

  // Add Game Events
  for i := 1; i < 18; i++ {
    var gameEvent GameEvent
    gameEvent.Season = 2014
    gameEvent.Week = i
    key := datastore.NewIncompleteKey(c, "GameEvent", tournyKey)
    if _, err := datastore.Put(c, key, &gameEvent); err != nil {
      panic(err.Error)
    }
  }

  utils.ServeJson(w, &t)
}
