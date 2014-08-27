package games

import (
  "appengine"
  "appengine/datastore"
  "net/http"

  "time"
  "github.com/go-martini/martini"
  "server/mikesnflpool/teams"
  "server/mikesnflpool/utils"
  "strconv"
)

type Tournament struct {
  Name          string            `json:"name"`
}

type GameEvent struct {
  TournamentKey *datastore.Key    `json:"tournament"`
  Season        int               `json:"season"`
  Week          int               `json:"week"`
}

type Game struct {
  Season          int             `json:"season"`
  Week            int             `json:"week"`
  Date            time.Time       `json:"date"`
  GameKey         *datastore.Key  `json:"gameKey" datastore:"-"` 
  Ended           bool            `json:"ended"`
 
  AwayTeamKey     *datastore.Key  `json:"awayTeamKey"`
  AwayTeamAbbr    string          `json:"awayTeamAbbr" datastore:"-"` 
  AwayTeam        teams.Team      `json:"awayTeam" datastore:"-"`
  AwayTeamScore   int             `json:"awayTeamScore"`
  AwayTeamSpread  float32         `json:"awayTeamSpread"`

  HomeTeamKey     *datastore.Key  `json:"homeTeamKey"`
  HomeTeamAbbr    string          `json:"homeTeamAbbr" datastore:"-"` 
  HomeTeam        teams.Team      `json:"homeTeam" datastore:"-"` 
  HomeTeamScore   int             `json:"homeTeamScore"`
  HomeTeamSpread  float32         `json:"homeTeamSpread"`
}

func GameHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  season,_ := strconv.Atoi(params["s"])
  week,_ := strconv.Atoi(r.URL.Query().Get("week"))

  // Get teams
  q := datastore.NewQuery("Team")
  var allTeams []teams.Team
  teamKeys, err := q.GetAll(c, &allTeams)
  if err != nil {
    panic(err.Error)
  }

  // Get games
  q = datastore.NewQuery("Game").
           Filter("Season =", season).
           Filter("Week =", week)
  var games []Game
  gameKeys, err := q.GetAll(c, &games)
  if err != nil {
    panic(err.Error)
  }

  // Associate team with game
  for i := range games {
    game := &games[i]
    game.GameKey = gameKeys[i]
    for j, teamKey := range teamKeys {
      if game.AwayTeamKey.Equal(teamKey) {
        game.AwayTeam = allTeams[j]
        continue
      }
      if game.HomeTeamKey.Equal(teamKey) {
        game.HomeTeam = allTeams[j]
        continue
      }
    }
  }

  utils.ServeJson(w, &games)
}

func AddOrUpdateGameHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var g Game
  if err := utils.ReadJson(r, &g); err != nil {
    panic(err.Error())
  }

  // Assign team keys
  g.AwayTeamKey = datastore.NewKey(c, "Team", g.AwayTeamAbbr, 0, nil)
  g.HomeTeamKey = datastore.NewKey(c, "Team", g.HomeTeamAbbr, 0, nil)

  // Get teams
  q := datastore.NewQuery("Team")
  var allTeams []teams.Team
  teamKeys, err := q.GetAll(c, &allTeams)
  if err != nil {
    panic(err.Error)
  }

  // Associate team with game
  for j, teamKey := range teamKeys {
    if g.AwayTeamKey.Equal(teamKey) {
      g.AwayTeam = allTeams[j]
      continue
    }
    if g.HomeTeamKey.Equal(teamKey) {
      g.HomeTeam = allTeams[j]
      continue
    }
  }

  // Check if existing game
  var existingGame Game
  err = datastore.Get(c, g.GameKey, &existingGame)
  isNew := err != nil 

  if isNew {
    // Add game
    key := datastore.NewIncompleteKey(c, "Game", nil)
    gameKey, err := datastore.Put(c, key, &g)
    if err != nil {
      panic(err.Error)
    }
    g.GameKey = gameKey
    utils.ServeJson(w, &g)
  } else {
    // Update existing game
    existingGame.AwayTeamScore = g.AwayTeamScore
    existingGame.HomeTeamScore = g.HomeTeamScore
    existingGame.AwayTeamSpread = g.AwayTeamSpread
    existingGame.HomeTeamSpread = g.HomeTeamSpread
    existingGame.Ended = g.Ended
    if _, err := datastore.Put(c, g.GameKey, &existingGame); err != nil {
      panic(err.Error)
    }
  }
}

func DeleteGameHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  gameKey, err := datastore.DecodeKey(params["g"])
  if err != nil {
    panic(err.Error)
  }

  // Get user pick keys for this game
  q := datastore.NewQuery("UserPick").Filter("GameKey =", gameKey).KeysOnly()
  keys, err := q.GetAll(c, nil)
  if err != nil {
    panic(err.Error)
  }

  // Delete keys
  keys = append(keys, gameKey)
  if err := datastore.DeleteMulti(c, keys); err != nil {
    panic(err.Error)
  }
}