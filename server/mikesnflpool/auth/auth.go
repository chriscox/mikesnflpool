package auth

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "server/mikesnflpool/utils"
)

type TestUser struct {
  Email     string    `json:"email"`
  Password  string    `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t TestUser
  if err := utils.ReadJson(r, &t); err != nil {
    panic(err.Error)
  }
  c.Infof("%v", t)
  key := datastore.NewIncompleteKey(c, "TestUser", nil)
  if _, err := datastore.Put(c, key, &t); err != nil {
    panic(err.Error)
  }
  utils.ServeJson(w, &t)
}
