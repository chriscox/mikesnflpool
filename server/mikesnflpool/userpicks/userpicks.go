package userpicks

import (
  "appengine"
  "appengine/datastore"
  "github.com/go-martini/martini"
  "net/http"
  "server/mikesnflpool/games"
  "server/mikesnflpool/teams"
  "server/mikesnflpool/tournaments"
  "server/mikesnflpool/utils"
  "strconv"
  "time"
  // "encoding/json"
)

type UserPick struct {
  Date              time.Time       `json:"date"`
  Game              games.Game      `json:"game" datastore:"-"`
  TeamKey           *datastore.Key  `json:"teamKey"`
  Team              teams.Team      `json:"team" datastore:"-"`
  GameKey           *datastore.Key  `json:"gameKey"`
  UserKey           *datastore.Key  `json:"userKey"`
  TournamentKey     *datastore.Key  `json:"tournamentKey" datastore:"-"`
  Season            int             `json:"season" datastore:"-"`
  Week              int             `json:"week" datastore:"-"`
}

// type WeeklyWins struct {
//   Week              int
//   Wins              int
// }

type WeeklyWins     map[string]int
type UserWins       map[string]WeeklyWins

type StatsMap struct {
  Stats              UserWins      `json:"stats"`
 }

// func NewWeeklyWins() *winsMap {
//   return &winsMap{userKey: "", weeklyWins: make(map[string]int)}
// }

func UserStatsHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  tournamentKey, err := datastore.DecodeKey(params["t"])
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
  // var userKeys []*datastore.Key
  // for i, _ := range tournamentUsers {
  //   userKeys = append(userKeys, tournamentUsers[i].UserKey)
  // }
  // var users = make([]User, len(userKeys))
  // if err := datastore.GetMulti(c, userKeys, users); err != nil {
  //   panic(err.Error)
  // }

  // Get all userpicks
  q = datastore.NewQuery("UserPick").Ancestor(tournamentKey)
  var allPicks []UserPick
  _, err = q.GetAll(c, &allPicks)
  if err != nil {
    panic(err.Error)
  }

  // Save game keys to associate game data later
  var gameKeys []*datastore.Key
  for _, pick := range allPicks {
    gameKeys = append(gameKeys, pick.GameKey)
  }

  // Get games with keys
  var games = make([]games.Game, len(allPicks))
  if err := datastore.GetMulti(c, gameKeys, games); err != nil {
    panic(err.Error)
  }

  // Associate game with picks
  for i, pick := range allPicks {
    for j, key := range gameKeys {
      if pick.GameKey.Equal(key) {
        allPicks[i].Game = games[j]
        break
      }
    }
  }

  // Calc wins
  var userWins = make(map[string]WeeklyWins)
  for _, u := range tournamentUsers {
    var weeklyWins = make(map[string]int)
    for _, p := range allPicks {
      if u.UserKey.Equal(p.UserKey) {
        if isSpreadWinner(p.Game, p) {
          weeklyWins[strconv.Itoa(p.Game.Week)] += 1
        }
      }
    }
    userWins[u.UserKey.Encode()] = weeklyWins
  }

  var stats StatsMap
  stats.Stats = userWins
  utils.ServeJson(w, &stats)
}

func isGameWinner(game games.Game, teamKey *datastore.Key) bool {
  if game.AwayTeamScore > game.HomeTeamScore {
    return teamKey.Equal(game.AwayTeamKey)
  } else if game.AwayTeamScore < game.HomeTeamScore {
    return teamKey.Equal(game.HomeTeamKey)
  } else {
    return false
  }
}

func isSpreadWinner(game games.Game, pick UserPick) bool {
  if game.Ended {
    if (float32(game.AwayTeamScore) - game.AwayTeamSpread) > (float32(game.HomeTeamScore) - game.HomeTeamSpread) {
      if (pick.TeamKey.Equal(game.AwayTeamKey)) {
        return true
      }
    } else if (float32(game.AwayTeamScore) - game.AwayTeamSpread) < (float32(game.HomeTeamScore) - game.HomeTeamSpread) {
      if (pick.TeamKey.Equal(game.HomeTeamKey)) {
        return true
      }
    } else {
      // If score minus spread equals other teams, then give losing
      // team the spread win. Teams must win spread + 1, not be equal.
      gameWinner := isGameWinner(game, game.AwayTeamKey)
      if !gameWinner {
        return pick.TeamKey.Equal(game.AwayTeamKey)
      } else {
        return pick.TeamKey.Equal(game.HomeTeamKey)
      }
    }
  }
  return false 
}

// TODO: Combine this and UserPickHandler
func AllUserPickHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  tournamentKey, err := datastore.DecodeKey(params["t"])
  if err != nil {
    panic(err.Error)
  }
  season,_ := strconv.Atoi(params["s"])
  week,_ := strconv.Atoi(r.URL.Query().Get("week"))

  // Get all teams
  var teams []teams.Team
  q := datastore.NewQuery("Team")
  teamKeys, err := q.GetAll(c, &teams)
  if err != nil {
    panic(err.Error)
  }

  // Get GameEvents ancestor for this season/week
  q = datastore.NewQuery("GameEvent").Ancestor(tournamentKey)
  var gameEvents []tournaments.GameEvent
  var gameEventKeyAncestor *datastore.Key
  gameEventKeys, err := q.GetAll(c, &gameEvents)
  if err != nil {
    panic(err.Error)
  }
  for i, e := range gameEvents {
    if e.Season == season && e.Week == week {
      gameEventKeyAncestor = gameEventKeys[i]
    }
  }

  // Get all pick for this game event ancestor
  var allPicks []UserPick
  q = datastore.NewQuery("UserPick").Ancestor(gameEventKeyAncestor)
  if _, err := q.GetAll(c, &allPicks); err != nil {
    panic(err.Error)
  }

  var gameKeys []*datastore.Key
  for iter1, pick := range allPicks {
    // Save game keys to associate game data later
    gameKeys = append(gameKeys, pick.GameKey)
    // Associate teams with picks
    teamKey := allPicks[iter1].TeamKey
    for iter2, t := range teams {
      if teamKey.Equal(teamKeys[iter2]) {
        allPicks[iter1].Team = t
        break
      }
    }
  }

  // Get games with keys
  var games = make([]games.Game, len(allPicks))
  if err := datastore.GetMulti(c, gameKeys, games); err != nil {
    panic(err.Error)
  }

  // Associate game with picks
  for i, pick := range allPicks {
    for j, key := range gameKeys {
      if pick.GameKey.Equal(key) {
        allPicks[i].Game = games[j]
        break
      }
    }
  }

  utils.ServeJson(w, &allPicks)
}
  
// TODO: This logic should match AllUserPickHandler, except it adds a filter to specific user.
func UserPickHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  // Get params
  tournamentKey, err := datastore.DecodeKey(params["t"])
  if err != nil {
    panic(err.Error)
  }
  userKey, err := datastore.DecodeKey(params["u"])
  if err != nil {
    panic(err.Error)
  }
  season,_ := strconv.Atoi(params["s"])
  week,_ := strconv.Atoi(r.URL.Query().Get("week"))

  // Get GameEvents ancestor for this season/week
  q := datastore.NewQuery("GameEvent").Ancestor(tournamentKey)
  var gameEvents []tournaments.GameEvent
  var gameEventKeyAncestor *datastore.Key
  gameEventKeys, err := q.GetAll(c, &gameEvents)
  if err != nil {
    panic(err.Error)
  }
  for i, e := range gameEvents {
    if e.Season == season && e.Week == week {
      gameEventKeyAncestor = gameEventKeys[i]
    }
  }

  // Get all pick for this game event ancestor
  var allPicks []UserPick
  var filteredPicks []UserPick
  q = datastore.NewQuery("UserPick").Ancestor(gameEventKeyAncestor)
  if _, err := q.GetAll(c, &allPicks); err != nil {
    panic(err.Error)
  }
  // Filter for this user
  for _, u := range allPicks {
    if u.UserKey.Equal(userKey) {
      filteredPicks = append(filteredPicks, u)
    }
  }

  utils.ServeJson(w, &filteredPicks)
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

  // Get existing userpick for this game
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
