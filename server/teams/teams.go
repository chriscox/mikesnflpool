package teams

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "server/utils"
)

type Team struct {
  Name      string  `json:"name"`
  NickName  string  `json:"abbr"`
}

func TeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Team").Limit(32)
  teams := make([]Team, 0, 32)
  if _, err := q.GetAll(c, &teams); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &teams)
}

func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t Team
  if err := utils.ReadJson(r, &t); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  key := datastore.NewKey(c, "Team", t.Name, 0, nil)
  if _, err := datastore.Put(c, key, &t); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &t)
}
