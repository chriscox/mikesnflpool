package auth

import (
  "appengine"
  "appengine/datastore"
  "net/http"
  "server/utils"
)

type TestUser struct {
  email  string
  password string
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var t TestUser
  if err := utils.ReadJson(r, &t); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  key := datastore.NewIncompleteKey(c, "TestUser", nil)
  if _, err := datastore.Put(c, key, &t); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &t)
}
