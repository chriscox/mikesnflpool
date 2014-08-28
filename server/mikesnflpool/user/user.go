package user

import (
  "appengine"
  "appengine/datastore"
  "bytes"
  "code.google.com/p/go.crypto/scrypt"
  "github.com/go-martini/martini"
  "net/http"
  "server/mikesnflpool/tournaments"
  "server/mikesnflpool/utils"
)

type User struct {
  FirstName       string            `json:"firstName"`
  LastName        string            `json:"lastName"`
  Email           string            `json:"email"`
  Password        string            `json:"password,omitempty" datastore:"-"`
  SecurePassword  []byte            `json:",omitempty"`
  UserKey         *datastore.Key    `json:"userKey" datastore:"-"`
  TournamentKey   *datastore.Key    `json:"tournamentKey" datastore:"-"`
  Admin           bool              `json:"admin,omitempty" datastore:"-"`
}

/*--- User Auth ---*/

func LoginHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var u User
  if err := utils.ReadJson(r, &u); err != nil {
    panic(err.Error)
  }

  // Get user
  var authUser User
  authUserKey := datastore.NewKey(c, "User", u.Email, 0, nil)
  err := datastore.Get(c, authUserKey, &authUser)
  if err != nil {
    panic(err.Error)
  }

  // Encrypt password and compare
  ctext, err := Encrypt(u.Password)
  if err != nil {
    panic(err.Error)
  }
  if !bytes.Equal(ctext, authUser.SecurePassword) {
    panic("Invalid login")
  }

  // Get tournament key
  q := datastore.NewQuery("TournamentUser").Filter("UserKey =", authUserKey)
  var tourneyUsers []tournaments.TournamentUser
  tourneyUserKeys, err := q.GetAll(c, &tourneyUsers)
  if err != nil {
    panic(err.Error)
  }
  if len(tourneyUserKeys) != 1 {
    panic(err.Error)
  }

  // IMPORTANT: clear password
  authUser.Password = "" 
  authUser.SecurePassword = nil

  // Send authenticated user
  authUser.UserKey = authUserKey
  authUser.Admin = tourneyUsers[0].Admin
  authUser.TournamentKey = tourneyUserKeys[0].Parent()
  utils.ServeJson(w, &authUser)
}

func UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var u User
  if err := utils.ReadJson(r, &u); err != nil {
    panic(err.Error)
  }

  // Create new user or quit if existing
  key := datastore.NewKey(c, "User", u.Email, 0, nil)
  if err := datastore.Get(c, key, &u); err == nil {
    w.WriteHeader(400)
    w.Write([]byte("A user with this email already exists."))
    return
  } 

  // Encrypt password
  ctext, err := Encrypt(u.Password)
  if err != nil {
    panic(err.Error)
  }
  u.SecurePassword = ctext

  // Save user
  userKey, err := datastore.Put(c, key, &u)
  if err != nil {
    panic(err.Error)
  }

  // Add Tournament User
  var t tournaments.TournamentUser
  t.UserKey = userKey
  key = datastore.NewIncompleteKey(c, "TournamentUser", u.TournamentKey)
  tourneyKey, err := datastore.Put(c, key, &t)
  if err != nil {
    panic(err.Error)
  }

  // IMPORTANT: clear passwords
  u.Password = "" 
  u.SecurePassword = nil

  // Send authenticated user
  u.UserKey = userKey
  u.TournamentKey = tourneyKey.Parent()
  utils.ServeJson(w, &u)
}

func UserHandler(parms martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  tournamentKey, err := datastore.DecodeKey(parms["t"])
  if err != nil {
    panic(err.Error)
  }

  // Get tournament users
  q := datastore.NewQuery("TournamentUser").Ancestor(tournamentKey)
  var tournamentUsers []tournaments.TournamentUser
  tournamentUserKeys, err := q.GetAll(c, &tournamentUsers)
  if err != nil {
    panic(err.Error)
  }

  // Build array of user keys and get users
  var userKeys []*datastore.Key
  for i, _ := range tournamentUsers {
    userKeys = append(userKeys, tournamentUsers[i].UserKey)
  }
  var users = make([]User, len(userKeys))
  if err := datastore.GetMulti(c, userKeys, users); err != nil {
    panic(err.Error)
  }

  // Send authenticated user array
  for i, k := range userKeys {
    users[i].UserKey = k
    // IMPORTANT: clear passwords
    users[i].Password = "" 
    users[i].SecurePassword = nil
    for j, t := range tournamentUserKeys {
      if k.Equal(tournamentUsers[j].UserKey) {
        users[i].TournamentKey = t.Parent()
        users[i].Admin = tournamentUsers[j].Admin
      }
    }
  }

  utils.ServeJson(w, &users)
}

func Encrypt(password string) ([]byte, error) {
  salt := []byte("%#?/*")
  return scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
}
