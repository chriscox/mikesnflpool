package user

import (
  "appengine"
  "appengine/datastore"
  // "github.com/go-martini/martini"
  "net/http"
  "server/mikesnflpool/games"
  "server/mikesnflpool/tournaments"
  "server/mikesnflpool/utils"
  // "strconv"
  "time"
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

type UserPick struct {
  Date              time.Time       `json:"date"`
  Game              games.Game      `json:"game" datastore:"-"`
  TeamKey           *datastore.Key  `json:"teamKey"`
  GameKey           *datastore.Key  `json:"gameKey"`
  UserKey           *datastore.Key  `json:"userKey"`
  TournamentKey     *datastore.Key  `json:"tournamentKey" datastore:"-"`
  Season            int             `json:"season" datastore:"-"`
  Week              int             `json:"week" datastore:"-"`
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
  if _, err = datastore.Put(c, key, &t); err != nil {
    panic(err.Error)
  }
  utils.ServeJson(w, &u)
}

/*--- User Picks ---*/
  
func UserPickHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var p UserPick
  if err := utils.ReadJson(r, &p); err != nil {
    panic(err.Error)
  }

  q := datastore.NewQuery("GameEvent").Ancestor(p.TournamentKey)
  iter := q.Run(c)
  for {
    var gameEvent tournaments.GameEvent
    gameEventKey, err := iter.Next(&gameEvent)
    if err == datastore.Done {
      break // No further entities match the query.
    }
    if err != nil {
      panic(err.Error)
    }
    if gameEvent.Season == p.Season && gameEvent.Week == p.Week {
      q = datastore.NewQuery("UserPick").Ancestor(gameEventKey)
      var picks []UserPick
      if _, err := q.GetAll(c, &picks); err != nil {
        panic(err.Error)
      }
      utils.ServeJson(w, &picks)
      break
    }
  }
}

func AddUserPickHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var p UserPick
  if err := utils.ReadJson(r, &p); err != nil {
    panic(err.Error)
  }

  // Set picked team
  game := p.Game
  p.GameKey = game.GameKey
  if game.AwayTeam.Selected {
    p.TeamKey = game.AwayTeamKey
  } else if game.HomeTeam.Selected {
    p.TeamKey = game.HomeTeamKey
  }

  // Get GameEvents with ancestor
  q := datastore.NewQuery("GameEvent").Ancestor(p.TournamentKey)
  var gameEvents []tournaments.GameEvent
  gameEventKeys, err := q.GetAll(c, &gameEvents)
  if err != nil {
    panic(err.Error)
  }

  // Update or add new pick
  q = datastore.NewQuery("UserPick").
          Filter("GameKey = ", p.GameKey).
          Filter("UserKey = ", p.UserKey)
  var existingPicks []UserPick
  existingPicksKeys, err := q.GetAll(c, &existingPicks)
  if err != nil {
    panic(err.Error)
  }

  if len(existingPicks) == 1 {
    // Update existing pick
    existingPicks[0].TeamKey = p.TeamKey
    if _, err := datastore.Put(c, existingPicksKeys[0], &existingPicks[0]); err != nil {
      panic(err.Error)
    }

  } else {
    // Save new UserPick
    for i, e := range gameEvents {
      if e.Season == game.Season && e.Week == game.Week {
        key := datastore.NewIncompleteKey(c, "UserPick", gameEventKeys[i])
        if _, err := datastore.Put(c, key, &p); err != nil {
          panic(err.Error)
        }
        break
      }
    }
  }

  utils.ServeJson(w, &p)
}
