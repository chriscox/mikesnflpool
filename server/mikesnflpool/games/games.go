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

type Game struct {
  Season        int             `json:"season"`
  Week          int             `json:"week"`
  Date          time.Time       `json:"date"`

  AwayTeamKey   *datastore.Key  `json:"awayTeam.key"`
  AwayTeamAbbr  string          `json:"awayTeamAbbr" datastore:"-"` 
  AwayTeam      teams.Team      `json:"awayTeam" datastore:"-"` 

  HomeTeamKey   *datastore.Key  `json:"homeTeam.key"`
  HomeTeamAbbr  string          `json:"homeTeamAbbr" datastore:"-"` 
  HomeTeam      teams.Team      `json:"homeTeam" datastore:"-"` 
}

func GameHandler(parms martini.Params, w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  season,_ := strconv.Atoi(parms["season"])
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
  if _, err := q.GetAll(c, &games); err != nil {
    panic(err.Error)
  }

  // Associate team with game
  for i := range games {
    awayTeamKey := games[i].AwayTeamKey.Encode()
    homeTeamKey := games[i].HomeTeamKey.Encode()
    for j, k := range teamKeys {
      teamKey := k.Encode()
      // Compare keys
      if awayTeamKey == teamKey {
        games[i].AwayTeam = allTeams[j]
        continue
      }
      if homeTeamKey == teamKey {
        games[i].HomeTeam = allTeams[j]
        continue
      }
    }
  }

  utils.ServeJson(w, &games)
}

func AddGameHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  var g Game
  if err := utils.ReadJson(r, &g); err != nil {
    panic(err.Error())
  }

  // Assign team keys
  g.AwayTeamKey = datastore.NewKey(c, "Team", g.AwayTeamAbbr, 0, nil)
  g.HomeTeamKey = datastore.NewKey(c, "Team", g.HomeTeamAbbr, 0, nil)

  // Add game
  key := datastore.NewIncompleteKey(c, "Game", nil)
  _, err := datastore.Put(c, key, &g)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  utils.ServeJson(w, &g)
}
