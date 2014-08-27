package user

import (
  "appengine"
  "appengine/datastore"
  "github.com/go-martini/martini"
  "net/http"
  "server/mikesnflpool/tournaments"
  "server/mikesnflpool/utils"
)

type User struct {
  FirstName       string            `json:"firstName"`
  LastName        string            `json:"lastName"`
  Email           string            `json:"email"`
  Password        string            `json:"password"`
  UserKey         *datastore.Key    `json:"userKey" datastore:"-"`
  TournamentKey   *datastore.Key    `json:"tournamentKey" datastore:"-"`
}

type AuthenticatedUser struct {
  FirstName       string            `json:"firstName"`
  LastName        string            `json:"lastName"`
  Email           string            `json:"email"`
  UserKey         *datastore.Key    `json:"userKey" datastore:"-"`
  TournamentKey   *datastore.Key    `json:"tournamentKey" datastore:"-"`
}

/*--- User Auth ---*/

func LoginHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var u User
  if err := utils.ReadJson(r, &u); err != nil {
    panic(err.Error)
  }

  // Check valid credentials
  q := datastore.NewQuery("User").
          Filter("Email =", u.Email).
          Filter("Password =", u.Password)
  var users []User
  userKeys, err := q.GetAll(c, &users)
  if err != nil {
    panic(err.Error)
  }
  if len(users) != 1 {
    panic("Invalid login")
  }

  q = datastore.NewQuery("TournamentUser").Filter("UserKey =", userKeys[0])
  var tourneyUsers []tournaments.TournamentUser
  tourneyUserKeys, err := q.GetAll(c, &tourneyUsers)
  if err != nil {
    panic(err.Error)
  }

  // Send authenticated user
  var user AuthenticatedUser
  user.FirstName = users[0].FirstName
  user.LastName = users[0].LastName
  user.Email = u.Email
  user.UserKey = userKeys[0]
  user.TournamentKey = tourneyUserKeys[0].Parent()
  utils.ServeJson(w, &user)
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
    panic("Existing User")
  } 
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

  // Send authenticated user
  var user AuthenticatedUser
  user.FirstName = u.FirstName
  user.LastName = u.LastName
  user.Email = u.Email
  user.UserKey = userKey
  user.TournamentKey = tourneyKey.Parent()
  utils.ServeJson(w, &user)
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
  _, err = q.GetAll(c, &tournamentUsers)
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
  var authenticatedUsers []AuthenticatedUser
  for j, u := range users {
    var user AuthenticatedUser
    user.FirstName = u.FirstName
    user.LastName = u.LastName
    user.Email = u.Email
    user.UserKey = userKeys[j]
    authenticatedUsers = append(authenticatedUsers, user)
  }

  utils.ServeJson(w, &authenticatedUsers)
}

