package teams

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "server/mikesnflpool/utils"
)

type Team struct {
  Abbr      string  `json:"abbr"`
  Name      string  `json:"name"`
  NickName  string  `json:"nickName"`
  Division  string  `json:"division"`
  Selected  bool    `json:"selected" datastore:"-"` 
}
  

func TeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Team")
  teams := make([]Team, 0, 32)
  if _, err := q.GetAll(c, &teams); err != nil {
    panic(err.Error)
  }
  utils.ServeJson(w, &teams)
}

func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t Team
  if err := utils.ReadJson(r, &t); err != nil {
    panic(err.Error)
  }

  key := datastore.NewKey(c, "Team", t.Abbr, 0, nil)
  if _, err := datastore.Put(c, key, &t); err != nil {
    panic(err.Error)
  }
  utils.ServeJson(w, &t)
}

// func GetAllTeams(r *http.Request) ([]Team, error) {
//   c := appengine.NewContext(r)
//   q := datastore.NewQuery("Team")
//   var teams []Team
//   if _, err := q.GetAll(c, &teams); err != nil {
//     return nil, err
//   }
//   return teams, nil
// }
